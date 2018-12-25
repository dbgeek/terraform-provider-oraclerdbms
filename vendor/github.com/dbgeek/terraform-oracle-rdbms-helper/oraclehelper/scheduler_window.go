package oraclehelper

import (
	"fmt"
	"log"
)

const (
	querySchedulerWindow = `
SELECT
	sw.owner,
    sw.window_name,
    sw.resource_plan,
    sw.schedule_type,
    sw.repeat_interval,
    TO_CHAR(sw.duration) AS duration,
    sw.enabled,
    sw.comments
FROM dba_scheduler_windows sw
WHERE sw.owner = UPPER(:1)
AND sw.window_name = UPPER(:2)
`
	execDropWindow = `
BEGIN
    DBMS_SCHEDULER.DROP_WINDOW (
		window_name => :1
	);
END;
`

	setAttribute = `
BEGIN
	DBMS_SCHEDULER.SET_ATTRIBUTE(
	name		=>	:1,
	attribute	=>	:2,
	value		=>	:3
);
END;
`
	disable = `
BEGIN
	DBMS_SCHEDULER.DISABLE(
		name	=>	:1
);
END;
`
	enable = `
BEGIN
	DBMS_SCHEDULER.ENABLE(
	name	=>	:1
);
END;
`
)

type (
	//ResourceSchedulerWindow ....
	ResourceSchedulerWindow struct {
		Owner          string
		WindowName     string
		ResourcePlan   string
		ScheduleType   string
		RepeatInterval string
		Duration       string
		Enabled        string
		StartDate      string
		Comments       string
		WindowPriority string
	}
	schedulerWindowService struct {
		client *Client
	}
)

func (s *schedulerWindowService) ReadSchedulerWindow(tf ResourceSchedulerWindow) (*ResourceSchedulerWindow, error) {
	log.Printf("[DEBUG] ReadSchedulerWindow windowname: %s\n", tf.WindowName)
	resourceSchedulerWindow := &ResourceSchedulerWindow{}

	err := s.client.DBClient.QueryRow(
		querySchedulerWindow,
		tf.Owner,
		tf.WindowName,
	).Scan(&resourceSchedulerWindow.Owner,
		&resourceSchedulerWindow.WindowName,
		&resourceSchedulerWindow.ResourcePlan,
		&resourceSchedulerWindow.ScheduleType,
		&resourceSchedulerWindow.RepeatInterval,
		&resourceSchedulerWindow.Duration,
		&resourceSchedulerWindow.Enabled,
		&resourceSchedulerWindow.Comments,
	)
	if err != nil {
		return nil, err
	}

	return resourceSchedulerWindow, nil
}

func (s *schedulerWindowService) CreateSchedulerWindow(tf ResourceSchedulerWindow) error {
	sqlCommand := fmt.Sprintf("BEGIN")
	sqlCommand += fmt.Sprintf(" DBMS_SCHEDULER.CREATE_WINDOW(")

	if tf.WindowName != "" {
		sqlCommand += fmt.Sprintf("window_name => '%s',", tf.WindowName)
	}
	if tf.ResourcePlan != "" {
		sqlCommand += fmt.Sprintf("resource_plan => '%s',", tf.ResourcePlan)
	}
	if tf.StartDate != "" {
		sqlCommand += fmt.Sprintf("start_date => %s,", tf.StartDate)
	}
	if tf.Duration != "" {
		sqlCommand += fmt.Sprintf("duration => '%s',", tf.Duration)
	}
	if tf.RepeatInterval != "" {
		sqlCommand += fmt.Sprintf("repeat_interval => '%s',", tf.RepeatInterval)
	}
	if tf.WindowPriority != "" {
		sqlCommand += fmt.Sprintf("window_priority => '%s',", tf.WindowPriority)
	}
	if tf.Comments != "" {
		sqlCommand += fmt.Sprintf("comments => '%s'", tf.Comments)
	}
	sqlCommand += fmt.Sprintf(");")
	sqlCommand += fmt.Sprintf("END;")

	log.Printf("[DEBUG] CreateSchedulerWindow sqlcommand: %s", sqlCommand)
	_, err := s.client.DBClient.Exec(sqlCommand)
	if err != nil {
		log.Printf("[ERROR] Create schedule Window failed with error: %v\n", err)
		return err
	}
	return nil
}

func (s *schedulerWindowService) DropSchedulerWindow(tf ResourceSchedulerWindow) error {
	_, err := s.client.DBClient.Exec(execDropWindow, tf.WindowName)
	if err != nil {
		log.Printf("[ERROR] drop schedule Window failed with error: %v\n", err)
		return err
	}
	return nil
}

func (s *schedulerWindowService) ModifySchedulerWindow(tf ResourceSchedulerWindow) error {
	name := fmt.Sprintf("%s.%s", tf.Owner, tf.WindowName)
	log.Printf("[DEBUG] ModifySchedulerWindow name: %s\n", name)
	_, err := s.client.DBClient.Exec(disable, name)

	if err != nil {
		log.Printf("[ERROR] ModifySchedulerWindow disable failed with error: %v\n", err)
	}

	if tf.ResourcePlan != "" {
		_, err := s.client.DBClient.Exec(setAttribute, name, "resource_plan", tf.ResourcePlan)
		if err != nil {
			log.Printf("[ERROR] ModifySchedulerWindow setAttribute failed with error: %v\n", err)
		}
	}
	if tf.StartDate != "" {
		_, err := s.client.DBClient.Exec(setAttribute, name, "start_date", tf.StartDate)
		if err != nil {
			log.Printf("[ERROR] ModifySchedulerWindow setAttribute failed with error: %v\n", err)
		}
	}
	if tf.Duration != "" {
		_, err := s.client.DBClient.Exec(setAttribute, name, "duration", tf.Duration)
		if err != nil {
			log.Printf("[ERROR] ModifySchedulerWindow setAttribute failed with error: %v\n", err)
		}
	}
	if tf.RepeatInterval != "" {
		_, err := s.client.DBClient.Exec(setAttribute, name, "repeat_interval", tf.RepeatInterval)
		if err != nil {
			log.Printf("[ERROR] ModifySchedulerWindow setAttribute failed with error: %v\n", err)
		}
	}
	if tf.WindowPriority != "" {
		_, err := s.client.DBClient.Exec(setAttribute, name, "window_priority", tf.WindowPriority)
		if err != nil {
			log.Printf("[ERROR] ModifySchedulerWindow setAttribute failed with error: %v\n", err)
		}
	}
	if tf.Comments != "" {
		_, err := s.client.DBClient.Exec(setAttribute, name, "comments", tf.Comments)
		if err != nil {
			log.Printf("[ERROR] ModifySchedulerWindow setAttribute failed with error: %v\n", err)
		}
	}
	_, err = s.client.DBClient.Exec(enable, name)
	if err != nil {
		log.Printf("[ERROR] ModifySchedulerWindow enabled failed with error: %v\n", err)
	}
	return nil
}
