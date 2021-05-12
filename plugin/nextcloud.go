package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Configuration struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
	Server   string `yaml:"server"`
	ApiPath  string `yaml:"api_path"`
}

func (c *Configuration) getConf() *Configuration {

	yamlFile, err := ioutil.ReadFile("/etc/check_mk/nextcloud.config.yml")
	//yamlFile, err := ioutil.ReadFile("./plugin/nextcloud.config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func printSegmentHeader(header string) {
	fmt.Println("<<<" + header + ">>>")
}

func printBasicInformation(info NextcloudInfo) {
	printSegmentHeader("nextcloud")
	fmt.Println(i2s(info.Ocs.Meta.Statuscode) + "|" + info.Ocs.Meta.Message)
	b, err := json.Marshal(info)
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	fmt.Println(string(b))
}

func printNextcloudSoftware(info NextcloudInfo) {
	printSegmentHeader("nextcloud_software")
	var system = info.Ocs.Data.Nextcloud.System
	fmt.Println(system.Version)
	fmt.Println(system.Freespace)
	var cpuLoad = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(system.Cpuload)), "|"), "[]")
	fmt.Println(cpuLoad)
	fmt.Println(i2s(system.MemTotal) + "|" + i2s(system.MemFree) + "|" + i2s(system.MemTotal-system.MemFree))
	fmt.Println(i2s(system.SwapTotal) + "|" + i2s(system.SwapFree) + "|" + i2s(system.SwapTotal-system.SwapFree))
}

func printNextcloudApps(info NextcloudInfo) {
	printSegmentHeader("nextcloud_apps")
	var apps = info.Ocs.Data.Nextcloud.System.Apps
	fmt.Println(i2s(apps.NumInstalled) + "|" + i2s(apps.NumUpdatesAvailable))
	b, err := json.Marshal(apps.AppUpdates)
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	fmt.Println(string(b))
}

func printNextcloudSystem(info NextcloudInfo) {
	printSegmentHeader("nextcloud_server")
	var opcache = info.Ocs.Data.Server.Php.Opcache
	var mem = opcache.MemoryUsage
	fmt.Println(i2s(mem.UsedMemory) + "|" + i2s(mem.FreeMemory) + "|" + i2s(mem.WastedMemory) + "|" + fmt.Sprintf("%f", mem.CurrentWastedPercentage))
}

func main() {
	var c Configuration
	c.getConf()

	info := getNextcloudJson(c)

	printBasicInformation(info)
	printNextcloudSoftware(info)
	printNextcloudApps(info)
	printNextcloudSystem(info)
}

func getNextcloudJson(c Configuration) NextcloudInfo {
	client := &http.Client{}
	req, err := http.NewRequest("GET", createUrl(c), nil)
	req.SetBasicAuth(c.Username, c.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)

	var data NextcloudInfo
	errUnm := json.Unmarshal(bodyText, &data)
	if errUnm != nil {
		log.Fatalf("Unmarshal: %v", errUnm)
	}
	return data
}

func createUrl(c Configuration) string {
	var url = c.Schema + "://" + c.Server + "/" + c.ApiPath
	return url
}

func i2s(i int) string {
	return strconv.Itoa(i)
}

type NextcloudInfo struct {
	Ocs struct {
		Meta struct {
			Status     string `json:"status"`
			Statuscode int    `json:"statuscode"`
			Message    string `json:"message"`
		} `json:"meta"`
		Data struct {
			Nextcloud struct {
				System struct {
					Version             string    `json:"version"`
					Theme               string    `json:"theme"`
					EnableAvatars       string    `json:"enable_avatars"`
					EnablePreviews      string    `json:"enable_previews"`
					MemcacheLocal       string    `json:"memcache.local"`
					MemcacheDistributed string    `json:"memcache.distributed"`
					FilelockingEnabled  string    `json:"filelocking.enabled"`
					MemcacheLocking     string    `json:"memcache.locking"`
					Debug               string    `json:"debug"`
					Freespace           int64     `json:"freespace"`
					Cpuload             []float64 `json:"cpuload"`
					MemTotal            int       `json:"mem_total"`
					MemFree             int       `json:"mem_free"`
					SwapTotal           int       `json:"swap_total"`
					SwapFree            int       `json:"swap_free"`
					Apps                struct {
						NumInstalled        int                    `json:"num_installed"`
						NumUpdatesAvailable int                    `json:"num_updates_available"`
						AppUpdates          map[string]interface{} `json:"app_updates"`
					} `json:"apps"`
				} `json:"system"`
				Storage struct {
					NumUsers         int `json:"num_users"`
					NumFiles         int `json:"num_files"`
					NumStorages      int `json:"num_storages"`
					NumStoragesLocal int `json:"num_storages_local"`
					NumStoragesHome  int `json:"num_storages_home"`
					NumStoragesOther int `json:"num_storages_other"`
				} `json:"storage"`
				Shares struct {
					NumShares               int `json:"num_shares"`
					NumSharesUser           int `json:"num_shares_user"`
					NumSharesGroups         int `json:"num_shares_groups"`
					NumSharesLink           int `json:"num_shares_link"`
					NumSharesMail           int `json:"num_shares_mail"`
					NumSharesRoom           int `json:"num_shares_room"`
					NumSharesLinkNoPassword int `json:"num_shares_link_no_password"`
					NumFedSharesSent        int `json:"num_fed_shares_sent"`
					NumFedSharesReceived    int `json:"num_fed_shares_received"`
				} `json:"shares"`
			} `json:"nextcloud"`
			Server struct {
				Webserver string `json:"webserver"`
				Php       struct {
					Version           string `json:"version"`
					MemoryLimit       int    `json:"memory_limit"`
					MaxExecutionTime  int    `json:"max_execution_time"`
					UploadMaxFilesize int    `json:"upload_max_filesize"`
					Opcache           struct {
						OpcacheEnabled    bool `json:"opcache_enabled"`
						CacheFull         bool `json:"cache_full"`
						RestartPending    bool `json:"restart_pending"`
						RestartInProgress bool `json:"restart_in_progress"`
						MemoryUsage       struct {
							UsedMemory              int     `json:"used_memory"`
							FreeMemory              int     `json:"free_memory"`
							WastedMemory            int     `json:"wasted_memory"`
							CurrentWastedPercentage float64 `json:"current_wasted_percentage"`
						} `json:"memory_usage"`
						InternedStringsUsage struct {
							BufferSize      int `json:"buffer_size"`
							UsedMemory      int `json:"used_memory"`
							FreeMemory      int `json:"free_memory"`
							NumberOfStrings int `json:"number_of_strings"`
						} `json:"interned_strings_usage"`
						OpcacheStatistics struct {
							NumCachedScripts   int     `json:"num_cached_scripts"`
							NumCachedKeys      int     `json:"num_cached_keys"`
							MaxCachedKeys      int     `json:"max_cached_keys"`
							Hits               int     `json:"hits"`
							StartTime          int     `json:"start_time"`
							LastRestartTime    int     `json:"last_restart_time"`
							OomRestarts        int     `json:"oom_restarts"`
							HashRestarts       int     `json:"hash_restarts"`
							ManualRestarts     int     `json:"manual_restarts"`
							Misses             int     `json:"misses"`
							BlacklistMisses    int     `json:"blacklist_misses"`
							BlacklistMissRatio int     `json:"blacklist_miss_ratio"`
							OpcacheHitRate     float64 `json:"opcache_hit_rate"`
						} `json:"opcache_statistics"`
					} `json:"opcache"`
					Apcu []interface{} `json:"apcu"`
				} `json:"php"`
				Database struct {
					Type    string `json:"type"`
					Version string `json:"version"`
					Size    int    `json:"size"`
				} `json:"database"`
			} `json:"server"`
			Activeusers struct {
				Last5Minutes int `json:"last5minutes"`
				Last1Hour    int `json:"last1hour"`
				Last24Hours  int `json:"last24hours"`
			} `json:"activeUsers"`
		} `json:"data"`
	} `json:"ocs"`
}
