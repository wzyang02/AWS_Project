package main        //everything belongs to a package 
import (
	 
	"fmt"
	"../pkg/mypkg"
)

//declare a varible can be accessed by any ffiles in same package

var myName string

//declare a function -- can be accessed by any files in same package

func myFunc(a string, b int, c []string) string {
	fmt.Println(a, b, "\n")
	for _,item := range c{
		fmt.Println(item)
	}
	return "Wonderful family"
}


//main function

func main(){
	var myName string
	myAge := 21
	myName = "William"

	myFamily := []string{"Haitian", "Yen", "Victoria"}
	message := myFunc(myName, myAge, myFamily) 
	fmt.Println(message)

	message = myPkg.GetAddress()
	fmt.Println(message)

	myPkg.GetPhone()

}



//run above code in terminal cmd directory
//go run ./main.go
//test