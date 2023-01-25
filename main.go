package main

import(
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/cobra"
)

// cria a estrutura do banco de dados User
type User struct {
	gorm.Model
	Username		string `gorm:"unique_index"`
	Name 			string
	Age				uint	
}

//Cria a conexão com o banco de dados sqlite3 com o nome test.db.
func main (){
	db, err := gorm.Open("sqlite3","test.db")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer db.Close()
	//Cria o banco de dados usando o que possui na struct User.
	db.AutoMigrate(&User{})
	//cria o rootCmd usando o package cobra
	var rootCmd = &cobra.Command{Use: "app"}
	//cria os comandos do cobra para utilizar no cmd usando flags.
	var createCmd = &cobra.Command{
		Use: "Create user",
		Short: "Create user",
		Run: func(cmd *cobra.Command, args []string){
			Username, _ := cmd.Flags().GetString("username")
			Name, _ := cmd.Flags().GetString("name")
			Age, _ := cmd.Flags().GetUint("age")
			db.Create(&User{Username: Username, Name: Name , Age: Age})
			fmt.Println("User created with successfully")
		},
	}
		//CreateCmd = cria os comandos usando as flags requerindo as flags respectivas
		createCmd.Flags().String("username", "", "username for user")
		createCmd.MarkFlagRequired("username")
		createCmd.Flags().String("name", "", "name for user")
		createCmd.MarkFlagRequired("name")
		createCmd.Flags().Uint("age", 0, "age for user")
		createCmd.MarkFlagRequired("age")
		rootCmd.AddCommand(createCmd)
		rootCmd.Execute()
	}