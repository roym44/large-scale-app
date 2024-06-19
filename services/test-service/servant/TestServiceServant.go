package TestServiceServant

var cacheMap map[string]string

func init() {
	cacheMap = make(map[string]string)
}

func HelloWorld() string {
    return "Hello World"
}

func HelloToUser(userName string) string {
    return "Hello " + userName
}

func Store(key string, value string) {
    cacheMap[key]=value
}

func Get(key string) (string,bool) {
    value,ok:=cacheMap[key]
    return value,ok
}
