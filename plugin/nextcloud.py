#!/usr/bin/env python3
import json
import os
from types import SimpleNamespace

import requests
from requests.auth import HTTPBasicAuth
import yaml

config_path = os.getenv('MK_CONFDIR') + '/nextcloud.config.yml'
config = yaml.safe_load(open(config_path))


def build_url():
    return config['schema'] + '://' + config['server'] + '/' + config['api_path']


def print_segment_header(name):
    print('<<<' + name + '>>>')


def get_info():
    return requests.get(build_url(), auth=HTTPBasicAuth(config['username'], config['password']))


def get_app_updates(response):
    return str(response.json()['ocs']['data']['nextcloud']['system']['apps']['app_updates'])


def print_information():
    print_segment_header('nextcloud')
    response = get_info()
    json_response = json.loads(response.content, object_hook=lambda d: SimpleNamespace(**d))
    ocs = json_response.ocs
    statuscode = ocs.meta.statuscode
    message = ocs.meta.message
    print(str(statuscode) + '|' + message)
    print(response.json())
    data = ocs.data
    system = data.nextcloud.system

    nextcloud_software(data,
                       system)

    nextcloud_apps(response, system)

    nextcloud_server(data)


def nextcloud_software(data,system ):
    print_segment_header('nextcloud_software')
    print(system.version)
    print(system.freespace)
    print('|'.join(str(s) for s in system.cpuload))
    print(str(system.mem_total) + '|' + str(system.mem_free) + '|' + str(system.mem_total - system.mem_free))
    print(str(system.swap_total) + '|' + str(system.swap_free) + '|' + str(system.swap_total - system.swap_free))


def nextcloud_apps(response, system):
    print_segment_header('nextcloud_apps')
    apps = system.apps
    print(str(apps.num_installed) + '|' + str(apps.num_updates_available) )
    print(get_app_updates(response))


def nextcloud_server(data):
    print_segment_header('nextcloud_server')
    opcache = data.server.php.opcache
    opcache_memory_usage = opcache.memory_usage
    print(str(opcache_memory_usage.used_memory) + '|' + str(opcache_memory_usage.free_memory) +
          '|' + str(opcache_memory_usage.wasted_memory) + '|' + str(opcache_memory_usage.current_wasted_percentage))


def print_json():
    print_segment_header('nextcloud')
    response = get_info()
    print(response.json())

print_information()
