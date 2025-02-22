package tui

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"git.qowevisa.me/Qowevisa/gotell/env"
	"git.qowevisa.me/Qowevisa/gotell/errors"
)

func GetIntPercentFromData(a, b int) int {
	return int(float64((a * 100)) / float64(b))
}

func SendMessageToConnectionEasy(msg *[]rune) (dataProcessHandler, dataT) {
	return SendMessageToConnection, dataT{rawP: msg}
}

func SendMessageToConnection(t *TUI, data dataT) error {
	if t.tlsConnection == nil {
		return errors.WrapErr("t.tlsConnection", errors.NOT_SET)
	}
	if data.rawP == nil {
		return errors.WrapErr("data.rawP", errors.NOT_SET)
	}
	message := string(*data.rawP)
	n, err := t.tlsConnection.Write([]byte(message))
	if err != nil {
		return errors.WrapErr("t.tlsConnection.Write", err)
	}
	log.Printf("Successfully wrote %d bytes to connection; Message: %s", n, message)
	return nil
}

// takes data from storage
func FE_ConnectTLS(t *TUI, data dataT) error {
	log.Printf("Start of FE_ConnectTLS")
	host, exist := t.storage[STORAGE_HOST_CONST]
	if !exist {
		errors.WrapErr("t.storage:host", errors.NOT_SET)
	}
	portStr, exist := t.storage[STORAGE_PORT_CONST]
	if !exist {
		errors.WrapErr("t.storage:host", errors.NOT_SET)
	}
	port, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		errors.WrapErr("port.strconv.ParseInt", err)
	}
	loadingFileName := env.ServerFullchainFileName
	cert, err := os.ReadFile(loadingFileName)
	if err != nil {
		errors.WrapErr("os.ReadFile", err)
	}
	log.Printf("Certificate %s loaded successfully!\n", loadingFileName)
	//
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(cert); !ok {
		errors.WrapErr("client: failed to parse root certificate", nil)
	}

	config := &tls.Config{
		RootCAs: roots,
	}
	if t.stateChannel == nil {
		return errors.WrapErr("t.stateChannel", errors.NOT_INIT)
	}
	t.stateChannel <- "TLS Connecting"
	conn, err := tls.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", host, int(port)),
		config,
	)
	if err != nil {
		t.stateChannel <- "TLS Connection error"
		return errors.WrapErr("tls.Dial", err)
	}
	t.stateChannel <- "TLS Established"
	t.tlsConnection = conn
	t.isConnected = true
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	t.tlsConnCloseData = closeData{
		wg:     wg,
		cancel: cancel,
	}
	go t.launchReadingMessagesFromConnection(ctx, wg)
	return nil
}

func CloseConnection(wg *sync.WaitGroup, cancel context.CancelFunc) {
	cancel()
	wg.Wait()
}

func AddToStorageEasy(key string, val *[]rune) (dataProcessHandler, dataT) {
	// that's why I create wrapper around it.
	//   try to understand that, dear viewer!
	return H_AddToStorage, dataT{rawP: val, op1: key}
}

func H_AddToStorage(t *TUI, data dataT) error {
	log.Printf("Debug: %#v", data)
	log.Printf("Adding to storage: %s = %s", data.op1, string(*data.rawP))
	if t.storage == nil {
		return errors.WrapErr("t.storage", errors.NOT_INIT)
	}
	t.storage[data.op1] = string(*data.rawP)
	return nil
}
