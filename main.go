package main

import (
	"bufio"
	"fmt"
	"html/template"
	
	"net/http"
	"os"
	"strings"
)

// Banner structure to store the banner (font) data
type Banner struct {
	Name  string
	Lines map[rune][]string
}
type data struct{
	Errr string
	Kalma string
}


// TemplateData structure for passing data to the HTML template
type TemplateData struct {
	Output string
}
type Error struct{
	Message string
}



func main() {
	// إعداد المسارات

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate", asciiArtHandler)

	


	// بدء الخادم على المنفذ 8080
	fmt.Println("star http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
func renderTemplate(w http.ResponseWriter, tmplPath string, data interface{}, statusCode int) {
	// Parse the template file
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the HTTP status code
	w.WriteHeader(statusCode)

	// Execute the template with the provided data
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}


// الصفحة الرئيسية التي تعرض النموذج
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// Parse the 404.html template
		ddddd:=data{
			Errr: "Error : Page Not Found",
			Kalma:"makayn walo",
		}
		renderTemplate(w,"templates/404.html",ddddd,404)
		return
	}

	// Parse the home page template (index.html)
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Execute the home page template
	tmpl.Execute(w, nil)
}

// التعامل مع طلب POST لتوليد ASCII Art
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {

	// استلام النص واختيار الخط
	inputText := r.FormValue("text")
	fontChoice := r.FormValue("font")


	// قراءة الخطوط (البنرات) من الملف المحدد
	banner, err := readBanner(fontChoice)
	if err != nil {
		http.Error(w, "Error reading font file", http.StatusInternalServerError)
		return
	}

	// تحويل النص إلى ASCII Art
	lines := strings.Split(inputText, "\n")
	output := convertToASCIIWithDynamicSpaces(lines, banner)
	if r.URL.Path != "/generate" {
		ddddd:=data{
			Errr: "Error : Page Not Found",
			Kalma:"makayn walo",
		}
		renderTemplate(w,"templates/404.html",ddddd,404)

		return
	}
	if len(inputText) == 0 {
		ddddd:=data{
			Errr: "Error : Page Not Found",
			Kalma:"makayn walo",
		}
		renderTemplate(w,"templates/404.html",ddddd,500)
				return
	}
	  for _,r:= range inputText{
		if r <=32 || r >= 126{
			ddddd:=data{
				Errr: "bad request",
				Kalma:"makayn walo",
			}
			renderTemplate(w,"templates/404.html",ddddd,400)
			return
		}
	  }
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		Output: output,
	}

	tmpl.Execute(w, data)
}

// قراءة ملف البنر وتحويله إلى خريطة
func readBanner(filename string) (*Banner, error) {
	file, err := os.Open("fonts/" + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	banner := &Banner{
		Name:  filename,
		Lines: make(map[rune][]string),
	}
	scanner := bufio.NewScanner(file)

	var currentChar rune
	charLines := []string{}
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// إذا وصلنا إلى سطر فارغ، نحفظ الحرف
		if line == "" {
			if len(charLines) > 0 {
				banner.Lines[currentChar] = charLines
			}
			charLines = []string{}
			lineCount = 0
			continue
		}

		// أول سطر يحدد الحرف
		if lineCount == 0 {
			currentChar = rune(len(banner.Lines) + 32) // الحروف تبدأ من ASCII 32 (المسافة)
		}

		// أضف السطر إلى الحرف الحالي
		charLines = append(charLines, line)
		lineCount++
	}

	// أضف آخر حرف
	if len(charLines) > 0 {
		banner.Lines[currentChar] = charLines
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return banner, nil
}


// تحويل النصوص متعددة الأسطر إلى ASCII مع الفراغات الديناميكية
func convertToASCIIWithDynamicSpaces(lines []string, banner *Banner) string {
	var result []string
	emptyLineCount := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			emptyLineCount++
			continue
		}
		if emptyLineCount > 0 {
			for i := 0; i < emptyLineCount; i++ {
				result = append(result, "")
			}
			emptyLineCount = 0
		}

		asciiLines := make([]string, 8) // Assuming the font has 8 rows of ASCII art
		for _, char := range line {
			// Check if the character exists in the banner
			if charLines, exists := banner.Lines[char]; exists {
				for i := 0; i < 8; i++ {
					asciiLines[i] += charLines[i]
				}
			} else {
				// If the character doesn't exist in the banner, add spaces
				for i := 0; i < 8; i++ {
					asciiLines[i] += " "
				}
			}
		}

		result = append(result, strings.Join(asciiLines, "\n"))
	}

	// Handle empty line padding at the end
	for i := 0; i < emptyLineCount; i++ {
		result = append(result, "")
	}

	return strings.Join(result, "\n")
}
