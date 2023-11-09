package encryption

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var AUTH_TOKEN_ECRYPTION_OFFSET byte = 15

//Encrypts a string with the format a:b into the file provided
func Encrypt(file string, data string) (error) {
	split_data := strings.Split(data, ":")
	d1 := []byte(split_data[0])

	for i := 0; i < len(d1); i++ {
		d1[i] += AUTH_TOKEN_ECRYPTION_OFFSET
	}

	d2 := []byte(split_data[1])

	for i := 0; i < len(d2); i++ {
		d2[i] += AUTH_TOKEN_ECRYPTION_OFFSET
	}
	
	encrypted_data := string(d1) + ":" + string(d2)
	
	wr_err := os.WriteFile(file, []byte(encrypted_data), 0644)
	
	if wr_err != nil {
		fmt.Println("File writing error", wr_err)
		return wr_err
	}
	
	return nil
}

func DecryptFromFile(file_name string) (string, string, error) {

	//Read in data from the file
	raw_data, r_err := ioutil.ReadFile(file_name)
	if r_err != nil {
		fmt.Println("File reading error", r_err)
		return "", "", r_err
	}

	//Split the two pieces of data
	split_data := strings.Split(string(raw_data), ":")
	data_1 := []byte(split_data[0])
	data_2 := []byte(split_data[1])

	//Decrypt the data
	for i := 0; i < len(data_1); i++ {
		data_1[i] -= AUTH_TOKEN_ECRYPTION_OFFSET
	}

	for i := 0; i < len(data_2); i++ {
		data_2[i] -= AUTH_TOKEN_ECRYPTION_OFFSET
	}

	return string(data_1), string(data_2), nil
}

//Grab the IP address needed to dconnect to the database from a text file
func GetIP() string {
	ip, r_err := ioutil.ReadFile("../encryption/db_ip.txt")
	if r_err != nil {
		fmt.Println("File reading error", r_err)
		return ""
	}

	return string(ip)
}