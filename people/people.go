package people

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)


//PeopleMap maps GitHub usernames to a human name where github doesnt have a listed User.Name
var PeopleMap map[string]string

type People struct {
	People []Person `json:"people"`
}

// User struct which contains a name
// a type and a list of social links
type Person struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

func init() {
	// Open our jsonFile
	usersjson, err := os.Open("users.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer usersjson.Close()

	byteValue, _ := ioutil.ReadAll(usersjson)

	// we initialize our Users array
	var people People

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &people)

	PeopleMap = make(map[string]string, len(people.People))
	for _, p := range people.People {
		PeopleMap[p.Username] = p.Name
	}

	fmt.Printf("Loaded peopleMap: %+v\n", PeopleMap)
}
