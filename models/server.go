package models

// DOGAPIURL is the url to fetch dog images from the dog api
// %s is the breed name
const DOGAPIURL = "https://dog.ceo/api/breed/%s/images/random"

// DOGAPIERRSTATUS is the status code for an error response from the dog api
// This is the status code returned when the dog api returns an error
const DOGAPIERRSTATUS = "error"
