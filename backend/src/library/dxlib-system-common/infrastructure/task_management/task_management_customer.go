package task_management

import (
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
)

func (tm *TaskManagement) DoCustomerCreate(l *log.DXLog, customerData utils.JSON) (id int64, err error) {
	id, err = tm.Customer.Insert(l, customerData)
	if err != nil {
		return 0, err
	}

	return id, nil
}
