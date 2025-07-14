package infrastructure

import (
	"github.com/donnyhardyanto/dxlib/configuration"
	"github.com/donnyhardyanto/dxlib/endpoint_rate_limiter"
	"github.com/donnyhardyanto/dxlib/redis"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/lib"
	"github.com/donnyhardyanto/dxlib_module/module/audit_log"
	"github.com/donnyhardyanto/dxlib_module/module/external_system"
	"github.com/donnyhardyanto/dxlib_module/module/general"
	"github.com/donnyhardyanto/dxlib_module/module/push_notification"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/arrears_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/configuration_settings"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/construction_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/master_data"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/upload_data"
	user_management_handler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/user_management/handler"
	"time"
)

func DoOnDefineSetVariables() (err error) {
	configRedis := *configuration.Manager.Configurations["redis"].Data
	configRedisRateLimit := configRedis["rate_limit"].(utils.JSON)
	redisNameIdRateLimit := redis.Manager.Redises[configRedisRateLimit["nameid"].(string)]

	configAPIEndpointRateLimiter := *configuration.Manager.Configurations["api_endpoint_rate_limiter"].Data
	v, ok := configAPIEndpointRateLimiter["default"]
	if !ok {
		return errors.Errorf("CONFIGURATION_DEFAULT_NOT_FOUND")
	}
	c := v.(utils.JSON)
	nameid, ok := c["nameid"].(string)
	if !ok {
		return errors.Errorf("CONFIGURATION_NAMEID_NOT_FOUND:%v", c)
	}
	maxAttempt, ok := c["max_attempts"].(int)
	if !ok {
		return errors.Errorf("CONFIGURATION_MAX_ATTEMPT_NOT_FOUND:%v", c)
	}
	timeWindowInMinutes, ok := c["time_window_in_minutes"].(int)
	if !ok {
		return errors.Errorf("CONFIGURATION_TIME_WINDOW_IN_MINUTES_NOT_FOUND:%v", c)
	}
	blockDurationInMinutes, ok := c["block_duration_in_minutes"].(int)
	if !ok {
		return errors.Errorf("CONFIGURATION_BLOCK_DURATION_IN_MINUTES_NOT_FOUND:%v", c)
	}
	endpoint_rate_limiter.Manager.Init(nameid,
		endpoint_rate_limiter.RateLimitConfig{
			MaxAttempts:   maxAttempt,
			TimeWindow:    time.Duration(timeWindowInMinutes) * time.Minute,
			BlockDuration: time.Duration(blockDurationInMinutes) * time.Minute,
		},
		redisNameIdRateLimit,
	)
	for k, v := range configAPIEndpointRateLimiter {
		if k == "default" {
			continue
		}
		c := v.(utils.JSON)
		nameid, ok := c["nameid"].(string)
		if !ok {
			return errors.Errorf("CONFIGURATION_NAMEID_NOT_FOUND:%v", c)
		}
		maxAttempt, ok := c["max_attempts"].(int)
		if !ok {
			return errors.Errorf("CONFIGURATION_MAX_ATTEMPT_NOT_FOUND:%v", c)
		}
		timeWindowInMinutes, ok := c["time_window_in_minutes"].(int)
		if !ok {
			return errors.Errorf("CONFIGURATION_TIME_WINDOW_IN_MINUTES_NOT_FOUND:%v", c)
		}
		blockDurationInMinutes, ok := c["block_duration_in_minutes"].(int)
		if !ok {
			return errors.Errorf("CONFIGURATION_BLOCK_DURATION_IN_MINUTES_NOT_FOUND:%v", c)
		}
		endpoint_rate_limiter.Manager.RegisterGroup(nameid, endpoint_rate_limiter.RateLimitConfig{
			MaxAttempts:   maxAttempt,
			TimeWindow:    time.Duration(timeWindowInMinutes) * time.Minute,
			BlockDuration: time.Duration(blockDurationInMinutes) * time.Minute,
		})
	}

	self.ModuleSelf.UserOrganizationMembershipType = user_management.UserOrganizationMembershipTypeSingleOrganizationPerUser
	self.ModuleSelf.OnAuthenticateUser = doOnAuthenticateUser
	self.ModuleSelf.OnCreateSessionObject = doOnCreateSessionObject

	configObjectStorage := *configuration.Manager.Configurations["object_storage"].Data

	configObjectStorageUserAvatarSource := configObjectStorage["user-avatar-source"].(utils.JSON)
	configObjectStorageUserAvatarSmall := configObjectStorage["user-avatar-small"].(utils.JSON)
	configObjectStorageUserAvatarMedium := configObjectStorage["user-avatar-medium"].(utils.JSON)

	configSecurity := *configuration.Manager.Configurations["security"].Data
	configSecurityImageUploader := configSecurity["image_uploader"].(utils.JSON)

	self.ModuleSelf.Avatar = lib.NewImageObjectStorage(configObjectStorageUserAvatarSource["nameid"].(string),
		configSecurityImageUploader["max_request_size"].(int64),
		configSecurityImageUploader["max_pixel_width"].(int64),
		configSecurityImageUploader["max_pixel_height"].(int64),
		configSecurityImageUploader["max_bytes_per_pixel"].(int64),
		configSecurityImageUploader["max_pixels"].(int64),
		map[string]lib.ProcessedImageObjectStorage{
			"small": {
				ObjectStorageNameId: configObjectStorageUserAvatarSmall["nameid"].(string),
				Width:               configObjectStorageUserAvatarSmall["file_image_width"].(int),
				Height:              configObjectStorageUserAvatarSmall["file_image_height"].(int),
			},
			"medium": {
				ObjectStorageNameId: configObjectStorageUserAvatarMedium["nameid"].(string),
				Width:               configObjectStorageUserAvatarMedium["file_image_width"].(int),
				Height:              configObjectStorageUserAvatarMedium["file_image_height"].(int),
			},
		})

	configObjectStorageAnnouncementPictureSource := configObjectStorage["announcement-picture-source"].(utils.JSON)
	configObjectStorageAnnouncementPictureSmall := configObjectStorage["announcement-picture-small"].(utils.JSON)
	configObjectStorageAnnouncementPictureMedium := configObjectStorage["announcement-picture-medium"].(utils.JSON)
	configObjectStorageAnnouncementPictureBig := configObjectStorage["announcement-picture-big"].(utils.JSON)

	general.ModuleGeneral.AnnouncementPicture = lib.NewImageObjectStorage(configObjectStorageAnnouncementPictureSource["nameid"].(string),
		configSecurityImageUploader["max_request_size"].(int64),
		configSecurityImageUploader["max_pixel_width"].(int64),
		configSecurityImageUploader["max_pixel_height"].(int64),
		configSecurityImageUploader["max_bytes_per_pixel"].(int64),
		configSecurityImageUploader["max_pixels"].(int64),
		map[string]lib.ProcessedImageObjectStorage{
			"small": {
				ObjectStorageNameId: configObjectStorageAnnouncementPictureSmall["nameid"].(string),
				Width:               configObjectStorageAnnouncementPictureSmall["file_image_width"].(int),
				Height:              configObjectStorageAnnouncementPictureSmall["file_image_height"].(int),
			},
			"medium": {
				ObjectStorageNameId: configObjectStorageAnnouncementPictureMedium["nameid"].(string),
				Width:               configObjectStorageAnnouncementPictureMedium["file_image_width"].(int),
				Height:              configObjectStorageAnnouncementPictureMedium["file_image_height"].(int),
			},
			"big": {
				ObjectStorageNameId: configObjectStorageAnnouncementPictureBig["nameid"].(string),
				Width:               configObjectStorageAnnouncementPictureBig["file_image_width"].(int),
				Height:              configObjectStorageAnnouncementPictureBig["file_image_height"].(int),
			},
		})

	configObjectStorageUserIdentityCardSource := configObjectStorage["user-identity-card-source"].(utils.JSON)
	configObjectStorageUserIdentityCardBig := configObjectStorage["user-identity-card-big"].(utils.JSON)

	user_management_handler.UserIdentityCard = lib.NewImageObjectStorage(configObjectStorageUserIdentityCardSource["nameid"].(string),
		configSecurityImageUploader["max_request_size"].(int64),
		configSecurityImageUploader["max_pixel_width"].(int64),
		configSecurityImageUploader["max_pixel_height"].(int64),
		configSecurityImageUploader["max_bytes_per_pixel"].(int64),
		configSecurityImageUploader["max_pixels"].(int64),
		map[string]lib.ProcessedImageObjectStorage{
			"big": {
				ObjectStorageNameId: configObjectStorageUserIdentityCardBig["nameid"].(string),
				Width:               configObjectStorageUserIdentityCardBig["file_image_width"].(int),
				Height:              configObjectStorageUserIdentityCardBig["file_image_height"].(int),
			},
		})

	user_management.ModuleUserManagement.UserOrganizationMembershipType = user_management.UserOrganizationMembershipTypeSingleOrganizationPerUser
	user_management.ModuleUserManagement.OnUserAfterCreate = user_management_handler.DoOnUserAfterCreate
	user_management.ModuleUserManagement.OnUserResetPassword = user_management_handler.DoOnUserResetPassword

	audit_log.ModuleAuditLog.Init(base.DatabaseNameIdAuditLog)
	configuration_settings.ModuleConfigurationSettings.Init(base.DatabaseNameIdConfig)
	external_system.ModuleExternalSystem.Init(base.DatabaseNameIdConfig)

	master_data.ModuleMasterData.Init(base.DatabaseNameIdTaskDispatcher)
	construction_management.ModuleConstructionManagement.Init(base.DatabaseNameIdTaskDispatcher)

	self.ModuleSelf.Init(base.DatabaseNameIdTaskDispatcher)

	general.ModuleGeneral.Init(base.DatabaseNameIdTaskDispatcher)

	user_management.ModuleUserManagement.Init(base.DatabaseNameIdTaskDispatcher)
	user_management.ModuleUserManagement.SessionRedis = redis.Manager.Redises["session"]
	user_management.ModuleUserManagement.PreKeyRedis = redis.Manager.Redises["prekey"]

	partner_management.ModulePartnerManagement.Init(base.DatabaseNameIdTaskDispatcher)

	push_notification.ModulePushNotification.FCM.Init(base.DatabaseNameIdTaskDispatcher)
	task_management.ModuleTaskManagement.Init(base.DatabaseNameIdTaskDispatcher)
	arrears_management.ModuleArrearsManagement.Init(base.DatabaseNameIdTaskDispatcher)
	upload_data.ModuleUploadData.Init(base.DatabaseNameIdTaskDispatcher)
	return nil
}
