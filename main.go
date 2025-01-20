package main

import (
	"fmt"
	"webcrawler/internal/utils"
)

func main() {
	input := `<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`
	fmt.Println("Running stuff")
	urls, err := utils.GetURLsFromHTML(input, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(urls)
}
