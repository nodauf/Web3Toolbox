package utilsDownloadSources

type GetSourceCodeEtherscanAPI struct {
	Message string `json:"message"`
	Result  []struct {
		Abi                  string `json:"ABI"`
		CompilerVersion      string `json:"CompilerVersion"`
		ConstructorArguments string `json:"ConstructorArguments"`
		ContractName         string `json:"ContractName"`
		EVMVersion           string `json:"EVMVersion"`
		Implementation       string `json:"Implementation"`
		Library              string `json:"Library"`
		LicenseType          string `json:"LicenseType"`
		OptimizationUsed     string `json:"OptimizationUsed"`
		Proxy                string `json:"Proxy"`
		Runs                 string `json:"Runs"`
		SourceCodeString     string `json:"SourceCode"`
		SourceCode           contractSourceCodeEtherscanAPI
		SwarmSource          string `json:"SwarmSource"`
	} `json:"result"`
	Status string `json:"status"`
}

type contractSourceCodeEtherscanAPI struct {
	Language string `json:"language"`
	Settings struct {
		Libraries struct{} `json:"libraries"`
		Optimizer struct {
			Enabled bool  `json:"enabled"`
			Runs    int64 `json:"runs"`
		} `json:"optimizer"`
		OutputSelection struct {
			_ struct {
				_ []string `json:"*"`
			} `json:"*"`
		} `json:"outputSelection"`
	} `json:"settings"`
	Sources map[string]struct {
		Content string `json:"content"`
	} `json:"sources"`
}
