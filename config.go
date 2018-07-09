package pygmaeus

import (
	"reflect"
	"strconv"
	"errors"
	"fmt"
	"flag"
	"os"
	"runtime"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var fileName string
var fileType FileType
var argOffset = 1

type FileType string

const (
	YmlExtension FileType = "yml"
	JsonExtension FileType = "json"
)

func SetFileType(t FileType)  {
	fileType = t
}

var isDebugMode bool

var configFlagSet = flag.NewFlagSet(FlagSetName, flag.ContinueOnError)

const FlagSetName = "pygmaeus-config"



func EnableDebug(enable bool) {
	isDebugMode = enable
}

func printIfDebug(format string, args ...interface{}) {
	if isDebugMode {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Printf("\n [%s:%d] ", fn, line)
		if len(args) > 0 {
			fmt.Printf(format, args...)
		} else {
			fmt.Print(format)
		}
	}
}

/*
	Prefix for environment variable taken by os.GetEnv;
	Set up by package function SetEnvPrefix
 */
var envPrefix = ""

/*
	Custom structure for taken flag from args.
	The aim of this tool take information that flag is exist at program
	args and then set up value of this flag in the structure field.
	The structure implement interface flag.Value
 */
type argFlag struct {
	reflectValue reflect.Value
	flagType     reflect.Kind
	set          bool
	value        string
}

func (fl *argFlag) Set(v string) error {
	switch fl.flagType {
	case reflect.String:
		fl.set = true
		fl.reflectValue.SetString(v)
	case reflect.Int:
		intVal, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		fl.set = true
		fl.reflectValue.SetInt(int64(intVal))
	case reflect.Int64:
		int64Val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		fl.set = true
		fl.reflectValue.SetInt(int64Val)
	case reflect.Float32:
		float32Val, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return err
		}
		fl.set = true
		fl.reflectValue.SetFloat(float32Val)
	case reflect.Float64:
		float64Val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		fl.set = true
		fl.reflectValue.SetFloat(float64Val)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}
		fl.set = true
		fl.reflectValue.SetBool(boolVal)
	default:
		return errors.New("unsupported type")
	}
	fl.value = v
	return nil
}

func (fl argFlag) String() string {
	return fl.value
}

func init() {
	fileName = "config"
	fileType = YmlExtension
}

/*
	Binding configuration to structure's argument in the next ordering:
	1. File configuration description.
	2. Env configuration description.
	3. Cli arguments configuration description.
 */

func Bind(v interface{}) {
	ReadFromFile(v)
	GetFromArgs(v)
}

func ReadFromFile(v interface{}) {
	var name string
	localConfig := fmt.Sprintf("%s_local.%s", fileName, fileType)
	if _, err := os.Stat(localConfig); os.IsNotExist(err){
		name = fmt.Sprintf("%s.%s", fileName, fileType)
	} else {
		name = localConfig
	}

	switch fileType {
	case YmlExtension:
		ReadFromYml(v,name)
	case JsonExtension:
		ReadFromJson(v,name)
	}
}

func ReadFromYml(v interface{}, filename string){
	printIfDebug("ReadFromYml: start  reading")
	defer printIfDebug("ReadFromYml: exit from function")
	dataByte, err := ioutil.ReadFile(filename)
	panicOnErr(err, "can't read config")
	err = yaml.Unmarshal(dataByte, v)
	panicOnErr(err,"can't unmarshal")
}

func ReadFromJson(v interface{},filename string)  {
	printIfDebug("ReadFromYml: start json reading")
	defer printIfDebug("ReadFromYml: exit from function")
	dataByte, err := ioutil.ReadFile(filename)
	panicOnErr(err,"can't read config")
	err = yaml.Unmarshal(dataByte, v)
	panicOnErr(err,"can't unmarshal")
}

func panicOnErr(err error, message string) {
	if err != nil {
		printIfDebug(message)
		panic(message)
	}
}

func GetFromArgs(v interface{}) {
	printIfDebug("GetFromArgs: start function")
	goRound(reflect.ValueOf(v), "")
	printIfDebug("GetFromArgs: start parsing")
	configFlagSet.Parse(os.Args[argOffset:])
}

func goRound(value reflect.Value, path string) {
	printIfDebug("goRound: start function")
	if value.Kind() == reflect.Ptr {
		printIfDebug("goRound: %s is pointer", path)
		value = value.Elem()
	}
	types := value.Type()
RootFor:
	for i := 0; i < value.NumField(); i++ {
		currentValue := value.Field(i)
		if !currentValue.IsValid() || !currentValue.CanSet() {
			continue RootFor
		}
		fieldType := types.Field(i)
		name := fieldType.Name
		fieldFlagName := path + name
		envVal, envValExist := os.LookupEnv(envPrefix + fieldFlagName)
		switch currentValue.Kind() {
		case reflect.Struct:
			printIfDebug("goRound: %s is structure", fieldFlagName)
			goRound(value.Field(i), fieldFlagName+".")
		case reflect.String:
			printIfDebug("goRound: %s flag is registered (string)", fieldFlagName)
			fl := &argFlag{flagType: reflect.String, reflectValue: currentValue}
			configFlagSet.Var(fl, fieldFlagName, fieldFlagName)
			if envValExist {
				currentValue.SetString(envVal)
			}
		case reflect.Int:
			printIfDebug("goRound: %s flag is registered (int)", fieldFlagName)
			fl := &argFlag{flagType: reflect.Int, reflectValue: currentValue}
			configFlagSet.Var(fl, fieldFlagName, fieldFlagName)
			if envValExist {
				intEnvVal, err := strconv.Atoi(envVal)
				if err == nil {
					currentValue.SetInt(int64(intEnvVal))
				} else {
					panic(err)
				}
			}
		case reflect.Int64:
			printIfDebug("goRound: %s flag is registered (int64)", fieldFlagName)
			fl := &argFlag{flagType: reflect.Int, reflectValue: currentValue}
			configFlagSet.Var(fl, fieldFlagName, fieldFlagName)
			if envValExist {
				intEnvVal, err := strconv.Atoi(envVal)
				if err == nil {
					currentValue.SetInt(int64(intEnvVal))
				} else {
					panic(err)
				}
			}
		case reflect.Float32:
			printIfDebug("goRound: %s flag is registered (float32)", fieldFlagName)
			fl := &argFlag{flagType: reflect.Float32, reflectValue: currentValue}
			configFlagSet.Var(fl, fieldFlagName, fieldFlagName)
			if envValExist {
				float32Val, err := strconv.ParseFloat(envVal, 32)
				if err == nil {
					currentValue.SetFloat(float64(float32Val))
				} else {
					panic(err)
				}
			}
		case reflect.Float64:
			printIfDebug("goRound: %s flag is registered (float64)", fieldFlagName)
			fl := &argFlag{flagType: reflect.Float64, reflectValue: currentValue}
			configFlagSet.Var(fl, fieldFlagName, fieldFlagName)
			if envValExist {
				float64Val, err := strconv.ParseFloat(envVal, 64)
				if err == nil {
					currentValue.SetFloat(float64Val)
				} else {
					panic(err)
				}
			}
		case reflect.Bool:
			printIfDebug("goRound: %s flag is registered (bool)", fieldFlagName)
			fl := &argFlag{flagType: reflect.Bool, reflectValue: currentValue}
			configFlagSet.Var(fl, fieldFlagName, fieldFlagName)
			if envValExist {
				boolVal, err := strconv.ParseBool(envVal)
				if err == nil {
					currentValue.SetBool(boolVal)
				} else {
					panic(err)
				}
			}
		}
	}
}
