package monkeylang

import (
	"fmt"
	"git.qowevisa.me/Qowevisa/gotell/tlep/gmyerr"
	"os"
	"time"
)

// don't care for errors for now
func canOpenDir(dirname string) bool {
	dir, err := os.Open(dirname)
	if err != nil {
		// TODO check for errors
		return false
	}
	stat, err := dir.Stat()
	if err != nil {
		// TODO check for errors
		return false
	}
	return stat.IsDir()
}

func GetDictionaryForUser(user string, getStat bool) (*Dictionary, error) {
	canI := canOpenDir(DictsDirName)
	if !canI {
		err := os.Mkdir(DictsDirName, 0755)
		if err != nil {
			return nil, gmyerr.WrapPrefix("os.Mkdir", err)
		}
	}
	exists, err := DoesDictExists(user)
	if err != nil {
		return nil, gmyerr.WrapPrefix("parser.DoWeCanLoadDict", err)
	}
	if !exists {
		if getStat {
			fmt.Println("Dictionary can't be find. Creating...")
		}
		befCreat := time.Now()
		dict, err := CreateNewDictionary()
		if err != nil {
			return nil, gmyerr.WrapPrefix("lang.CreateNewDictionary", err)
		}
		if getStat {
			fmt.Println("Creation time is ", time.Now().Sub(befCreat))
			fmt.Println("Saving to file...")
		}
		befSave := time.Now()
		err = SaveToFile(*dict, user)
		if err != nil {
			return nil, gmyerr.WrapPrefix("parser.SaveToFile", err)
		}
		if getStat {
			fmt.Println("Saving time is ", time.Now().Sub(befSave))
		}
		return dict, nil
	}
	if getStat {
		fmt.Println("Dictionary can be find. Loading...")
	}
	befLoad := time.Now()
	dict, err := LoadFromFile(user)
	if err != nil {
		return nil, gmyerr.WrapPrefix("parser.LoadFromFile", err)
	}
	if getStat {
		fmt.Println("Loading time is ", time.Now().Sub(befLoad))
	}
	return dict, nil
}
