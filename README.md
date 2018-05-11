![Build Status](https://travis-ci.org/adtechpotok/go-common.svg?branch=master)

## Общие либы гошных демонов

##DB writer 
It sends data to database. If data base is unavailable, after 2 tries it will write into files in configurated directory. When connection is back, it will resend data to db in same query orders.
Datatable must contains field_server_id, to support maxinstanced daemons
##Example
```$xslt
	db,_  := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
         		//Mysql.User,
         		//Mysql.Password,
         		//Mysql.Host,
         		//Mysql.Port,
         		//Table.Name,
         		"utf-8"))
         		
	writer := dbwriter.New(dbwriter.WriteConfig{
		Db:                db,
		Log:               logrus.New(),
		FilePath:          mainConfig.DataWriter.FilePath,
		ServerId:          mainConfig.ServerId,
		TickTimeMs:        500, // it will send aggregated data to db every half second
		MaxConnectTimeSec: 2*60, // time, which for mysql connection is alive
	})
	//than 
    data := struct {
    	id      int       `gorm:column:id`
    	name    string    `gorm:column:name`
    	created time.Time `gorm:column:created`
    }{1, "Name", time.Now()}
    writer.Append(data)    
```
*do not create additional writer instance without need
*those queries are not injections save