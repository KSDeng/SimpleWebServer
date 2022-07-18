package config

type Configurations struct {
	WebServer WebServerConfigs `json:"webserver"`
	TCPServer TCPServerConfigs `json:"tcpserver"`
	Database  DatabaseConfigs  `json:"database"`
	Redis     RedisConfigs     `json:"redis"`
}

type WebServerConfigs struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type TCPServerConfigs struct {
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
}

type DatabaseConfigs struct {
	DB_User          string `json:"db_user"`
	DB_Password      string `json:"db_password"`
	DB_Name          string `json:"db_name"`
	DB_Net           string `json:"db_net"`
	DB_Addr          string `json:"db_addr"`
	DB_Driver        string `json:"db_driver"`
	DB_Max_Conn      int    `json:"db_max_conn"`
	DB_Max_Idle_Conn int    `json:"db_max_idle_conn"`
}

type RedisConfigs struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}
