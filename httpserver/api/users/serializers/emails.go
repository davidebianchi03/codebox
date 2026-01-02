package serializers

type TestEmailResponseSerializer struct {
	Success     bool   `json:"success"`
	Description string `json:"description"`
}

func LoadTestEmailResponseSerializer(
	success bool, description string,
) TestEmailResponseSerializer {
	return TestEmailResponseSerializer{
		Success:     success,
		Description: description,
	}
}
