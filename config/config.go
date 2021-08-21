package config

import (
	"fmt"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/view_models"
	"joranvest/repository"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type configService struct {
	applicationUserRepository repository.ApplicationUserRepository
}

func (db *configService) setupAdministrator() (bool, error) {
	var superAdmin = models.ApplicationUser{
		FirstName: "Administrator",
		LastName:  "",
		Username:  "sys_admin",
		Password:  "joranvest",
		Address:   "",
		Phone:     "",
		IsAdmin:   true,
		IsDefault: true,
		Email:     "admin@joranvest.com",
	}
	newRecord, err := db.applicationUserRepository.Insert(superAdmin)
	if err != nil {
		return false, err
	} else {
		return newRecord.IsActive, nil
	}
}

//-- Setup and Open Database Connection
func SetupDatabaseConnection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env file")
	}

	//dbDial := os.Getenv("DB_DIAL")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	fmt.Println("Open Connection..")
	fmt.Println("Connecting...")
	if err != nil {
		panic("Failed to connect database.")
	} else {
		fmt.Println("Connected.")
	}

	//-- This function to generate model to database table
	db.AutoMigrate(
		&models.Entity{},
		&models.ApplicationUser{},
		&models.Role{},
		&models.RoleMember{},

		&models.Emiten{},
		&models.Sector{},
		&models.Tag{},
		&models.ArticleCategory{},
		&models.Article{},
		&models.ArticleTag{},
		&models.Membership{},
		&models.Filemaster{},

		&models.TechnicalAnalysis{},
		&models.FundamentalAnalysis{},
		&models.FundamentalAnalysisTag{},

		&models.Order{},
	)

	viewList := make(map[string]map[string]string)

	var vw_order = entity_view_models.EntityMembershipView{}
	viewList[vw_order.TableName()] = vw_order.Migration()

	var vw_emiten = entity_view_models.EntityEmitenView{}
	viewList[vw_emiten.TableName()] = vw_emiten.Migration()

	var vw_sector = entity_view_models.EntitySectorView{}
	viewList[vw_sector.TableName()] = vw_sector.Migration()

	var vw_filemaster = entity_view_models.EntityFilemasterView{}
	viewList[vw_filemaster.TableName()] = vw_filemaster.Migration()

	var vw_tag = entity_view_models.EntityTagView{}
	viewList[vw_tag.TableName()] = vw_tag.Migration()

	var vw_technical_analysis = entity_view_models.EntityTechnicalAnalysisView{}
	viewList[vw_technical_analysis.TableName()] = vw_technical_analysis.Migration()

	var vw_fundamental_analysis = entity_view_models.EntityFundamentalAnalysisView{}
	viewList[vw_fundamental_analysis.TableName()] = vw_fundamental_analysis.Migration()

	var vw_fundamental_analysis_tag = entity_view_models.EntityFundamentalAnalysisTagView{}
	viewList[vw_fundamental_analysis_tag.TableName()] = vw_fundamental_analysis_tag.Migration()

	if len(viewList) > 0 {
		for _, detail := range viewList {
			db.Exec(fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s", detail["view_name"], detail["query"]))
		}
	}

	var entityRecord = models.Entity{}
	if err := db.Where("Name = ?", "Joranvest").First(&entityRecord).Error; err != nil {
		fmt.Println("Creating Administrator Started")
		entityRecord.Id = uuid.New().String()
		entityRecord = models.Entity{
			Name:        "Joranvest",
			Description: "Joranvest",
		}
		entityRecord.Id = uuid.New().String()
		db.Create(&entityRecord)
		fmt.Println("New Entity has been created.")
	}

	var superAdmin = models.ApplicationUser{}
	if err := db.Where("Username = ?", os.Getenv("DEFAULT_ADMINISTRATOR_USERNAME")).First(&superAdmin).Error; err != nil {
		fmt.Println("Creating Administrator Started")
		superAdmin.Id = uuid.New().String()
		superAdmin = models.ApplicationUser{
			FirstName: "System",
			LastName:  "Administrator",
			Username:  os.Getenv("DEFAULT_ADMINISTRATOR_USERNAME"),
			Address:   "",
			Phone:     "",
			IsAdmin:   true,
			Email:     "admin@joranvest.com",
		}
		superAdmin.Id = uuid.New().String()
		superAdmin.CreatedBy = superAdmin.Id
		superAdmin.Password = helper.HashAndSalt([]byte("joranvest"))
		superAdmin.EntityId = entityRecord.Id
		db.Create(&superAdmin)
		fmt.Println("Finished and Enjoy.")
	}
	return db
}

//-- Close Database Connection
func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to close connection")
	}
	dbSQL.Close()
}
