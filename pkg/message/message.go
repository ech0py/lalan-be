package message

/*
	Konstanta untuk pesan respons API.

Berisi pesan umum dan spesifik untuk berbagai operasi.
*/
const (
	// Pesan umum
	MsgSuccess              = "Success."
	MsgBadRequest           = "Bad request."
	MsgUnauthorized         = "Unauthorized."
	MsgNotAllowed           = "Method not allowed."
	MsgInternalServerError  = "Internal server error."
	MsgUserIDRequired       = "User ID is required."
	MsgFailedToHashPassword = "Failed to hash password."

	// Pesan autentikasi hoster
	MsgHosterCreatedSuccess     = "Hoster created successfully."
	MsgHosterCreatedFailed      = "Failed to create hoster."
	MsgHosterLoginSuccess       = "Hoster logged in successfully."
	MsgHosterLoginFailed        = "Hoster login failed."
	MsgHosterInvalidEmail       = "Invalid email format."
	MsgHosterInvalidCredentials = "Invalid email or password."
	MsgHosterEmailExists        = "A hoster with this email already exists."
	MsgHosterWeakPassword       = "The provided password is too weak."
	MsgHosterNotFound           = "Hoster not found."
	MsgHosterFetched            = "Hoster data retrieved successfully."

	// Pesan autentikasi customer
	MsgCustomerCreatedSuccess     = "Customer created successfully."
	MsgCustomerCreatedFailed      = "Failed to create customer."
	MsgCustomerLoginSuccess       = "Customer logged in successfully."
	MsgCustomerLoginFailed        = "Customer login failed."
	MsgCustomerInvalidEmail       = "Invalid email format."
	MsgCustomerInvalidCredentials = "Invalid email or password."
	MsgCustomerEmailExists        = "A customer with this email already exists."
	MsgCustomerWeakPassword       = "The provided password is too weak."
	MsgCustomerNotFound           = "Customer not found."
	MsgCustomerFetched            = "Customer data retrieved successfully."

	// Pesan kategori
	MsgCategoryCreatedSuccess = "Category created successfully."
	MsgCategoryUpdatedSuccess = "Category updated successfully."
	MsgCategoryDeletedSuccess = "Category deleted successfully."
	MsgCategoryNameExists     = "Category name already exists."
	MsgCategoryNotFound       = "Category not found."
	MsgCategoryNameRequired   = "Category name is required."
	MsgCategoryNameTooLong    = "Category name must not exceed 255 characters."
	MsgCategoryIDRequired     = "Category ID is required."

	// Pesan item
	MsgItemCreatedSuccess     = "Item created successfully."
	MsgItemUpdatedSuccess     = "Item updated successfully."
	MsgItemDeletedSuccess     = "Item deleted successfully."
	MsgItemNameExists         = "Item name already exists."
	MsgItemNotFound           = "Item not found."
	MsgItemNameRequired       = "Item name is required."
	MsgItemNameTooLong        = "Item name must not exceed 255 characters."
	MsgItemIDRequired         = "Item ID is required."
	MsgItemStockInvalid       = "Item stock cannot be negative."
	MsgItemPricePerDayInvalid = "Item price per day cannot be negative."
	MsgItemDepositInvalid     = "Item deposit cannot be negative."

	// Pesan terms and conditions
	MsgTermAndConditionsCreatedSuccess      = "Terms and conditions created successfully."
	MsgTermAndConditionsUpdatedSuccess      = "Terms and conditions updated successfully."
	MsgTermAndConditionsDeletedSuccess      = "Terms and conditions deleted successfully."
	MsgTermAndConditionsAlreadyExists       = "Terms and conditions already exist."
	MsgTermAndConditionsNotFound            = "Terms and conditions not found."
	MsgTermAndConditionsDescriptionRequired = "Terms and conditions description is required."
	MsgTermAndConditionsDescriptionTooShort = "Terms and conditions description must be at least 255 characters."
	MsgTermAndConditionsIDRequired          = "Terms and conditions ID is required."
	MsgTermAndConditionsStockInvalid        = "Terms and conditions stock cannot be negative."
	MsgTermAndConditionsPricePerDayInvalid  = "Terms and conditions price per day cannot be negative."
	MsgTermAndConditionsDepositInvalid      = "Terms and conditions deposit cannot be negative."
)
