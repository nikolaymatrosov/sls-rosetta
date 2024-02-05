package main

//
//import (
//	"context"
//	"encoding/json"
//	"log"
//	"os"
//
//	ycsdk "github.com/yandex-cloud/go-sdk"
//	"github.com/yandex-cloud/go-sdk/iamkey"
//)
//
//func main() {
//	iamToken := os.Getenv("YC_IAM_TOKEN")
//	var config ycsdk.Config
//	if iamToken != "" {
//		config = ycsdk.Config{
//			Credentials: ycsdk.NewIAMTokenCredentials(iamToken),
//		}
//	} else {
//		saJsonPath := os.Getenv("YC_SA_JSON_PATH")
//		if saJsonPath == "" {
//			log.Fatal("YC_SA_JSON_PATH or YC_IAM_TOKEN must be set")
//		}
//
//		saJsonFile, err := os.ReadFile(saJsonPath)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		var saKey iamkey.Key
//
//		err = json.Unmarshal(saJsonFile, &saKey)
//
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		key, err := ycsdk.ServiceAccountKey(&saKey)
//		if err != nil {
//			return
//		}
//		config = ycsdk.Config{
//			Credentials: key,
//		}
//	}
//
//	sdk, err := ycsdk.Build(context.Background(), config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Fragment from Pushkin's "The Stationmaster"
//	text := "Кто не проклинал станционных смотрителей, кто с ними не бранивался? Кто, в минуту гнева, не требовал от них роковой книги, дабы вписать в оную свою бесполезную жалобу на притеснение, грубость и неисправность? Кто не почитает их извергами человеческого рода, равными покойным подьячим или по крайней мере муромским разбойникам? Будем, однако, справедливы, постараемся войти в их положение и, может быть, станем судить о них гораздо снисходительнее."
//
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Println(p)
//}
