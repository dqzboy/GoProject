package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	HarborDN       = "www.dqzboy.com"     // 填入访问域名或IP
	HarborAddress  = "https://" + HarborDN // Harbor的URL地址(需要带前缀http或者https)
	HarborUser     = "admin"              // 登录Harbor的用户
	HarborPassword = "Harbor123456"    // 登录Harbor的用户密码
)

func main() {
	projects := []string{"public"} // 项目列表，可以添加多个项目

	for _, project := range projects {
		imageNames, err := getRepositoryNames(project)
		if err != nil {
			fmt。Printf("Error getting repository names for project %s: %v\n", project, err)
			continue
		}

		for _, imageName := range imageNames {
			imageTags, err := getTagsForRepository(project, imageName)
			if err != nil {
				fmt。Printf("Error getting tags for image %s in project %s: %v\n", imageName, project, err)
				continue
			}

			for _, tag := range imageTags {
				fmt。Printf("%s/%s:%s\n", HarborDN, imageName, tag)
			}
		}
	}
}

func getRepositoryNames(project string) ([]string， error) {
	url := fmt。Sprintf("%s/api/v2.0/projects/%s/repositories", HarborAddress, project)
	response, err := sendRequest(url)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	if err := json。Unmarshal(response, &data); err != nil {
		return nil, err
	}

	var names []string
	for _, repo := range data {
		name, ok := repo["name"]。(string)
		if ok {
			names = append(names, name)
		}
	}

	return names, nil
}

func getTagsForRepository(project, imageName string) ([]string， error) {
	url := fmt。Sprintf("%s/v2/%s/tags/list", HarborAddress, imageName)
	response, err := sendRequest(url)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json。Unmarshal(response, &data); err != nil {
		return nil, err
	}

	tagsInterface, ok := data["tags"]。([]interface{})
	if !ok {
		return nil, fmt。Errorf("Tags not found in response")
	}

	var tags []string
	for _, tag := range tagsInterface {
		tagStr, ok := tag。(string)
		if ok {
			tags = append(tags, tagStr)
		}
	}

	return tags, nil
}

func sendRequest(url string) ([]byte， error) {
	req, err := http。NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req。SetBasicAuth(HarborUser, HarborPassword)
	req。Header。Set("Content-Type", "application/json")

	client := &http。Client{}
	resp, err := client。Do(req)
	if err != nil {
		return nil, err
	}
	defer resp。Body。Close()

	body, err := ioutil。ReadAll(resp。Body)
	if err != nil {
		return nil, err
	}

	if resp。StatusCode != http。StatusOK {
		return nil, fmt。Errorf("Request failed with status code %d", resp。StatusCode)
	}

	return body, nil
}
