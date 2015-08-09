package main

import (
	"os"
	"os/exec"
)

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var defaultCert = map[string]string{
	// AOSP
	"android_media":    "B7:9D:F4:A8:2E:90:B5:7E:A7:65:25:AB:70:37:AB:23:8A:42:F5:D3",
	"android_platform": "27:19:6E:38:6B:87:5E:76:AD:F7:00:E7:EA:84:E4:C6:EE:E3:3D:FA",
	"android_release":  "61:ED:37:7E:85:D3:86:A8:DF:EE:6B:86:4B:D8:5B:0B:FA:A5:AF:81",
	"android_shared":   "5B:36:8C:FF:2D:A2:68:69:96:BC:95:EA:C1:90:EA:A4:F5:63:0F:E5",
}

func readCert(certFile string) map[string]string {
	cert := make(map[string]string)
	for k, v := range defaultCert {
		cert[k] = v
	}

	// load cert file into memory
	if certFile != "" {
		data, err := ioutil.ReadFile(certFile)
		if err != nil {
			log.Print(err)
			return nil
		}

		certs := strings.Split(string(data), "\n")
		matcher := regexp.MustCompile(`([^\s]+)\s+([^\s]+)\s+([^\s]+)$`)
		commentMatcher := regexp.MustCompile(`^\s*#`)
		for _, line := range certs {
			res := commentMatcher.FindStringSubmatch(line)
			if res != nil {
				continue
			}
			res = matcher.FindStringSubmatch(line)
			if res == nil {
				continue
			}
			key := fmt.Sprintf("%s_%s", res[1], res[2])
			cert[key] = res[3]
		}
	}

	return cert
}

func getSha1HashInFile(filename string) (result string) {
	result = ""
	shaMatcher := regexp.MustCompile(`\s+SHA1: (.+)`)
	// TODO: find an alternative way to get sha1 hash
	// call keytool -printcert -file <filename> |
	cmd := exec.Command("keytool", "-printcert", "-file", filename)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}
	o, err := ioutil.ReadAll(stdout)
	output := strings.Split(string(o), "\n")
	if err := cmd.Wait(); err != nil {
		log.Println(err)
		return
	}
	for _, line := range output {
		res := shaMatcher.FindStringSubmatch(line)
		if res != nil {
			return res[1]
		}
	}
	return
}

func getSha1Hash(apkFile string) (apkSha1Hash string) {
	// open apk file
	r, err := zip.OpenReader(apkFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Close()

	matcher := regexp.MustCompile(`META-INF/.+.RSA$`)

	for _, f := range r.File {
		res := matcher.FindStringSubmatch(f.Name)
		if res == nil {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			log.Println(err)
			continue
		}
		defer rc.Close()

		rsaData, err := ioutil.ReadAll(rc)
		if err != nil {
			log.Println(err)
			continue
		}

		fileName := os.TempDir() + "cert"
		os.Remove(fileName)
		f, err := os.Create(fileName)
		if err != nil {
			log.Println(err)
			continue
		}
		defer f.Close()

		if _, err := f.Write(rsaData); err != nil {
			log.Println(err)
			continue
		}
		// remember to delete file
		defer os.Remove(fileName)

		return getSha1HashInFile(fileName)
	}

	fmt.Println("There's no RSA file")
	return
}

// Checkcert return false if vendor/certType incorrect
func Checkcert(apkFile, certFile, vendor, certType string) (result bool) {
	cert := readCert(certFile)
	if cert == nil {
		return false
	}

	apkSha1Hash := getSha1Hash(apkFile)
	if apkSha1Hash == "" {
		return false
	}

	// get cert
	apkVendor := ""
	apkCertType := ""
	for k, v := range cert {
		if v == apkSha1Hash {
			s := strings.Split(k, "_")
			apkVendor = s[0]
			apkCertType = s[1]
		}
	}

	if apkVendor == "" || apkCertType == "" {
		fmt.Printf("%s is using \x1b[41mUNKNOWN\x1b[m key: %s\n", apkFile, apkSha1Hash)
	} else {
		fmt.Printf("%s is using \x1b[32m%s\x1b[m's \x1b[33m%s\x1b[m key\n", apkFile, apkVendor, apkCertType)
	}

	if vendor != "" && apkVendor != vendor {
		return false
	}

	if certType != "" && apkCertType != certType {
		return false
	}

	return true
}
