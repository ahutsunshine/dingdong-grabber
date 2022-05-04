package config

// Config 字段含义详情请参考config.yaml
type Config struct {
	Cookie                string   `yaml:"cookie"`
	Strategy              int      `yaml:"strategy"`
	CronJobs              []string `yaml:"cron_jobs"`
	BaseThreadSize        int      `yaml:"base_thread_size"`
	SubmitOrderThreadSize int      `yaml:"submit_order_thread_size"`
	MinSleepMillis        int      `yaml:"min_sleep_millis"`
	MaxSleepMillis        int      `yaml:"max_sleep_millis"`
	PushToken             string   `yaml:"push_token"`
	Play                  bool     `yaml:"play"`
}

type Conf struct {
	Name   string
	Config *Config `yaml:"dingdong"`
}
