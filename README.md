# go_follina

go run main.go -h
Encode b64 : 


$string= {PWSH SCRIPT}
[System.Convert]::ToBase64String([System.Text.Encoding]::UNICODE.GetBytes($string))
