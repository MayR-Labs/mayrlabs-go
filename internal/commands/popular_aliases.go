package commands

// GetPopularAliases returns a list of popular aliases
func GetPopularAliases() []Alias {
	return []Alias{
		// Git
		{Name: "gs", Command: "git status"},
		{Name: "ga", Command: "git add ."},
		{Name: "gc", Command: "git commit -m"},
		{Name: "gp", Command: "git push"},
		{Name: "gl", Command: "git pull"},
		{Name: "gco", Command: "git checkout"},
		{Name: "gd", Command: "git diff"},
		{Name: "glg", Command: "git log --oneline --graph --decorate"},
		{Name: "gcb", Command: "git checkout -b"},
		{Name: "gbr", Command: "git branch"},
		{Name: "grs", Command: "git restore --staged ."},
		{Name: "gsta", Command: "git stash push"},
		{Name: "gstp", Command: "git stash pop"},

		// System
		{Name: "ll", Command: "ls -la"},
		{Name: "cls", Command: "clear"},
		{Name: "update", Command: "sudo apt update && sudo apt upgrade -y"},
		{Name: "..", Command: "cd .."},
		{Name: "...", Command: "cd ../.."},
		{Name: "h", Command: "history"},
		{Name: "j", Command: "jobs"},

		// Docker
		{Name: "dps", Command: "docker ps"},
		{Name: "dpsa", Command: "docker ps -a"},
		{Name: "dco", Command: "docker-compose"},
		{Name: "dcup", Command: "docker-compose up -d"},
		{Name: "dcdown", Command: "docker-compose down"},
		{Name: "dlogs", Command: "docker-compose logs -f"},
		{Name: "dstop", Command: "docker stop $(docker ps -a -q)"},
		{Name: "drm", Command: "docker rm $(docker ps -a -q)"},
		{Name: "drmi", Command: "docker rmi $(docker images -q)"},

		// Go
		{Name: "gt", Command: "go test ./..."},
		{Name: "gr", Command: "go run ."},
		{Name: "gb", Command: "go build ."},
		{Name: "gmt", Command: "go mod tidy"},
		{Name: "gqv", Command: "go version"},

		// Node/NPM
		{Name: "ni", Command: "npm install"},
		{Name: "ns", Command: "npm start"},
		{Name: "nb", Command: "npm run build"},
		{Name: "nt", Command: "npm test"},
		{Name: "nd", Command: "npm run dev"},

		// Yarn
		{Name: "yi", Command: "yarn install"},
		{Name: "ys", Command: "yarn start"},
		{Name: "yb", Command: "yarn build"},
		{Name: "yd", Command: "yarn dev"},
	}
}
