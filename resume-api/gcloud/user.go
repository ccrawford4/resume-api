package gcloud

// Key = name, value = content of the resume
type ResumeContent map[string]string

// TODO: Actually fetch a users resumes from the Google Cloud Storage bucket and return a list of resume files
func fetchAndDownloadResumes(userId string) []string {
	return ResumeContent{
		""
	}
}
