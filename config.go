package pygmaeus

import (
	"reflect"
	"strconv"
	"errors"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"flag"
	"os"
)

var defaultFile string
var defaultType FileType

type FileType string

const (
	FILE_YML  FileType = ".yml"
	FILE_JSON FileType = ".json"
)

var fileName = "config"
/*
	Prefix for environment variable taken by os.GetEnv;
	Set up by package function SetEnvPrefix
 */
var envPrefix = "" //

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
		fl.reflectValue.SetString(fl.String())
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
	defaultFile = "config.yml"
}

func Bind(v interface{}) {
	ReadFromFile(v)
}

func ReadFromFile(v interface{}) {
	dataByte, err := ioutil.ReadFile(defaultFile)
	if err != nil {
		panic("can't read config")
	}
	err = yaml.Unmarshal(dataByte, v)
	if err != nil {
		panic("can't unmarshal")
	}
}

func GetFromArgs(v interface{}) {
	goRound(reflect.ValueOf(v), "")
}

func goRound(value reflect.Value, path string) {
	if value.Kind() == reflect.Ptr {
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
		fmt.Println(fieldFlagName, currentValue.Kind())
		envVal, envValExist := os.LookupEnv(fieldFlagName)
		switch currentValue.Kind() {
		case reflect.Struct:
			goRound(value.Field(i), name+".")
		case reflect.String:
			fl := &argFlag{flagType: reflect.String, reflectValue: currentValue}
			flag.Var(fl, fieldFlagName, fieldFlagName)
			if envValExist {
				currentValue.SetString(envPrefix + envVal)
			}
		case reflect.Int:
			fl := &argFlag{flagType: reflect.Int, reflectValue: currentValue}
			flag.Var(fl, fieldFlagName, fieldFlagName)
			if envValExist {
				intEnvVal, err := strconv.Atoi(envVal)
				if err == nil {
					currentValue.SetInt(int64(intEnvVal))
				} else {
					panic(err)
				}
			}
		case reflect.Int64:
			fl := &argFlag{flagType: reflect.Int, reflectValue: currentValue}
			flag.Var(fl, fieldFlagName, fieldFlagName)
			if envValExist {
				intEnvVal, err := strconv.Atoi(envVal)
				if err == nil {
					currentValue.SetInt(int64(intEnvVal))
				} else {
					panic(err)
				}
			}
		case reflect.Float32:
			fl := &argFlag{flagType: reflect.Float32, reflectValue: currentValue}
			flag.Var(fl, fieldFlagName, fieldFlagName)
			if envValExist {
				float32Val, err := strconv.ParseFloat(envVal, 32 )
				if err == nil {
					currentValue.SetFloat(float64(float32Val))
				} else {
					panic(err)
				}
			}
		case reflect.Float64:
			fl := &argFlag{flagType: reflect.Float64, reflectValue: currentValue}
			flag.Var(fl, fieldFlagName, fieldFlagName)
			if envValExist {
				float64Val, err := strconv.ParseFloat(envVal, 64 )
				if err == nil {
					currentValue.SetFloat(float64Val)
				} else {
					panic(err)
				}
			}
		}
	}
	if len(path) == 0 {
		flag.Parse()
	}
}

func SetDefaultFileType(t FileType) {
	defaultType = t
}

func SetEnvPrefix(s string) {
	envPrefix = s
}
