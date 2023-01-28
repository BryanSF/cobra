package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/cobra"
)

// cria a estrutura do banco de dados User
type User struct {
	gorm.Model
	Username string `gorm:"unique_index"`
	Name     string
	Age      uint
	// Chave estrangeira da Struct Company.
	CompanyID    uint    `gorm:"ForeignKey:CompanyRefer"`
	CompanyRefer Company `gorm:"ForeignKey:CompanyID;AssociationForeignkey:ID"`
}

//Cria a estrutura do banco de dados da tabela Company
type Company struct {
	gorm.Model
	Name string `gorm:"unique_index"`
}

//Cria a conexão com o banco de dados sqlite3 com o nome test.db.
func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	//Cria o banco de dados usando o que possui na struct User.
	db.AutoMigrate(&User{}, &Company{})
	//cria o rootCmd usando o package cobra
	var rootCmd = &cobra.Command{Use: "app"}
	//cria os comandos do cobra para utilizar no cmd usando flags.
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create user",
		Run: func(cmd *cobra.Command, args []string) {
			Username, _ := cmd.Flags().GetString("username")
			Name, _ := cmd.Flags().GetString("name")
			Age, _ := cmd.Flags().GetUint("age")
			CompanyID, _ := cmd.Flags().GetUint("companyID")
			nameCompany, _ := cmd.Flags().GetString("nameCompany")
			tx := db.Begin() // tx = transacao onde comeca um push request
			if err := tx.Create(&Company{Name: nameCompany}).Error; err != nil { // verifica se todos os dados estao ok
				tx.Rollback() // se nao estiverem ele reseta tudo e nao passa
				fmt.Println("Error ao criar copanhia", err)
				return// retorna o erro
			}
			if err := tx.Create(&User{Username: Username, Name: Name, Age: Age, CompanyID: CompanyID}).Error; err != nil { // verifica se todos os dados estao ok
				tx.Rollback() // se nao estiverem ele reseta tudo e nao passa
				fmt.Println("Error ao criar usuario.", err)
				return // retorna o erro
			}
			tx.Commit() //se tudo estiver ok ele lanca no banco o que foi pedido
			fmt.Println("User and company created with successfully")
		},
	}
	var searchCmd = &cobra.Command{ //cria uma variavel de procura 
		Use:   "search",
		Short: "Search in db table",
		Run: func(cmd *cobra.Command, args []string) {
			Username, _ := cmd.Flags().GetString("username")
			nameCompany, _ := cmd.Flags().GetString("nameCompany")
			var user User //Variavel de usuario colocando a struct 
			db.Where("username = ?", Username).First(&user) //Procura o username solicitado no banco de dados users
			fmt.Println("User:", user)
			var company Company //Variavel de company com struct company
			db.Where("nameCompany = ?", nameCompany).First(&company) //Procura o username solicitado no banco de dados companies
			fmt.Println("nameCompany:", company)
		},
	}
	//CreateCmd = cria os comandos usando as flags requerindo as flags respectivas
	createCmd.Flags().String("username", "", "username for user")
	createCmd.MarkFlagRequired("username")
	createCmd.Flags().String("name", "", "name for user")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().Uint("age", 0, "age for user")
	createCmd.MarkFlagRequired("age")
	createCmd.Flags().Uint("companyID", 0, "companyID for user")
	createCmd.MarkFlagRequired("companyID")
	createCmd.Flags().String("nameCompany", "", "nameCompany for user")
	createCmd.MarkFlagRequired("nameCompany")
	//comandos cobra
	searchCmd.Flags().String("username", "", "search username in user")
	searchCmd.MarkFlagRequired("username")
	searchCmd.Flags().String("nameCompany", "", "search nameCompany in company")
	searchCmd.MarkFlagRequired("nameCompany")
	//adicionado searchCmd no rootCmd
	rootCmd.AddCommand(createCmd, searchCmd)
	rootCmd.Execute()
}
