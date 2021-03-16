from .agent_based_api.v1 import *


def discover_linux_nextcloud(section):
    yield Service()


def check_linux_nextcloud(section):
    response = section[0][0].split('|')
    response_code = int(response[0])
    response_message = response[1]

    if response_code != 200:
        yield Result(state=State.CRIT, summary="Response="+str(response_code) + " Message=" + response_message)
        return
    yield Result(state=State.OK, summary="Response="+str(response_code) + " Message=" + response_message)


def to_mb(bytes):
    return float(bytes) / 1000 / 1000


def warn_level(max):
    return int(float(max) * 0.75)


def crit_level(max):
    return int(float(max) * 0.9)


def check_linux_nextcloud_software(section):
    version = section[0][0]
    freespace = section[1][0]
    cpuload = [float(i) for i in section[2][0].split('|')]
    mem = [to_mb(int(i)) for i in section[3][0].split('|')]
    swap = [to_mb(int(i)) for i in section[4][0].split('|')]

    yield Metric(
        "mem",
        mem[2],
        levels=(warn_level(mem[0]),crit_level(mem[0])),
        boundaries=(0,mem[0]))

    yield Metric(
        "swap",
        swap[2],
        levels=(warn_level(swap[0]),crit_level(swap[0])),
        boundaries=(0,swap[0]))

    error = False

    if mem[2] > warn_level(mem[0]):
        yield Result(state=State.WARN, summary="Used Memory > 75%. Max|Free|Used=" + str(mem))
        error = True
    elif mem[2] > crit_level(mem[0]):
        yield Result(state=State.CRIT, summary="Used Memory > 90%. Max|Free|Used=" + str(mem))
        error = True

    if swap[2] > warn_level(swap[0]):
        yield Result(state=State.WARN, summary="Used SWAP > 75%. Max|Free|Used=" + str(swap))
        error = True
    elif  swap[2] > crit_level(swap[0]):
        yield Result(state=State.CRIT, summary="Used SWAP > 90%. Max|Free|Used=" + str(swap))
        error = True

    if not error:
        yield Result(state=State.OK, summary="Version="+version)


def check_linux_nextcloud_server(section):
    opcache = section[0][0].split('|')
    max = to_mb(int(opcache[0]) + int(opcache[1]))

    yield Metric(
        "opcache",
        to_mb(int(opcache[0])),
        levels=(warn_level(max), crit_level(max)),
        boundaries=(0, max))
    yield Result(state=State.OK, summary="OPCache (used|free|wasted|wasted%)=" + str(opcache))


def check_linux_nextcloud_apps(section):
    response = section[0][0].split('|')
    installed_apps = int(response[0])
    updates = int(response[1])

    if updates > 0:
        yield Result(state=State.CRIT, summary="Update Available =" + str(section[1]))
    else:
        yield Result(state=State.OK, summary="Installed Apps = " + str(installed_apps))



register.check_plugin(
    name="nextcloud",
    service_name="Nextcloud",
    discovery_function=discover_linux_nextcloud,
    check_function=check_linux_nextcloud,
)


register.check_plugin(
    name="nextcloud_software",
    service_name="Nextcloud Software",
    discovery_function=discover_linux_nextcloud,
    check_function=check_linux_nextcloud_software,
)

register.check_plugin(
    name="nextcloud_apps",
    service_name="Nextcloud Apps",
    discovery_function=discover_linux_nextcloud,
    check_function=check_linux_nextcloud_apps,
)

register.check_plugin(
    name="nextcloud_server",
    service_name="Nextcloud Server",
    discovery_function=discover_linux_nextcloud,
    check_function=check_linux_nextcloud_server,
)
