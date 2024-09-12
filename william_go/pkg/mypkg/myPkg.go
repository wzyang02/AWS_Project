package myPkg        //everything belongs to a package
import (
	 
	"fmt"
)



var myAddress string   //first letter lower case is private, can only be access in same package 
 
 

func GetAddress() string {
	myAddress := "\n4331 Wedgewood Dr Copley Oh 44321"
	fmt.Println(myAddress, "\n")
	
	return "Good Place to live"
}

func GetPhone() {
	myPhone := 3304009725
	fmt.Println(myPhone, "\n")

	return 
}
 
 