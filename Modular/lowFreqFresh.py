from __future__ import unicode_literals

import json
import logging
import os
# import random
import time
from datetime import datetime, timedelta
from math import floor, log
from xml.dom.minidom import parse
# from redis.exceptions import RedisError
import numpy as np
import pandas as pd
import requests

from lib.Config.config import LOW_FREQ, MACHINE_CONFIG, get_conf_file
from lib.Interface.mysqlConnector import (select_from_single_machine,
                                          select_single_machine_with_raw)
from lib.Interface.redis_api import publish_message
from lib.Interface.redisConnector import RedisConnector
from lib.utils import date_format, reduce_program_name, get_tag_value

from lib.Interface.opcConnector import get_tool_compensation_to_opc

# 找到软件系统根目录
# ROOT_FOLDER = "IMS" # 根文件夹
# DEPAND_FOLDER = "99_Layer_Depend" # 依赖文件夹
# path_abso = pathlib.Path().absolute() # 当前文件绝对路径
# path_part = path_abso.parts # 文件路径截取
# path_root = pathlib.Path().joinpath(*path_part[0:path_part.index(ROOT_FOLDER)+1]) # 找到根文件夹路径
# path_depend = path_root.joinpath(DEPAND_FOLDER) # 找到配置文件绝对路径
# sys.path.append(str(path_depend))
# APIS_CONFIG = os.path.join(str(path_depend), "API_config.xml")


# path = os.path.abspath(os.path.dirname(__file__))
# FILE = os.path.join(path, "backend.xlsx")

OLD_DATA = ""
TOOL = {}
# H6000_WORKPIECE_FINISHED_COUNT = 0
# H8000_WORKPIECE_FINISHED_COUNT = 0
# NHM6300_WORKPIECE_FINISHED_COUNT = 0
MACHINE_NAME = ""
# FILE = get_conf_file("backend.xlsx")
machines_config = get_conf_file(MACHINE_CONFIG)
domTree = parse(machines_config)
rootNode = domTree.documentElement
machine = rootNode.getElementsByTagName("Machine")

status_mapping = pd.read_csv(os.path.join('conf', 'mock_machine_status_mapping.csv'))
MACHINE_STATUS = status_mapping.iloc[0].to_dict()
MACHINE_STATUS_INV = {v: k for k, v in MACHINE_STATUS.items()}

df_workpiece = pd.read_csv(os.path.join('conf', 'mock_workpiece.csv'))
# 构建 workpiece 字典
workpiece = {
    str(workpiece_id): {"product": product, "product_id": product_id, "finished_number": finished_number,
                        "product_remain": product_remain}
    for workpiece_id, product, product_id, finished_number, product_remain in zip(
        df_workpiece["id"], df_workpiece["product"], df_workpiece["product_id"], df_workpiece["finished_number"],
        df_workpiece["product_remain"]
    )
}


# machine_specifications = pd.read_csv(os.path.join('conf', 'mock_machine_specifications.csv'))
# machine_basic_information = machine_specifications.iloc[0, :6].to_dict()
# # print(machine_basic_information)
#

# machine_constraints = machine_specifications.iloc[0, 6:9].to_dict()
# # print(machine_constraints)
#

# machine_tool_magazine = machine_specifications.iloc[0, 9:14].to_dict()
# # print(machine_tool_magazine)
#

# machine_cool = machine_specifications.iloc[0, 14:17].to_dict()
# # print(machine_cool)
#

# machine_precision = machine_specifications.iloc[0, 17:22].to_dict()
# # print(machine_precision)
#

# machine_power = {"主轴功率曲线": []}
# machine_power_curve = machine_specifications["主轴功率曲线"].tolist()[1:]
# point = [i for i in range(1, len(machine_power_curve) + 1)]
# machine_power["主轴功率曲线"] = [point, machine_power_curve]
# # print(machine_power)
#

# machine_torque = {"主轴扭矩曲线": []}
# machine_torque_curve = machine_specifications["主轴扭矩曲线"].tolist()[1:]
# point = [i for i in range(1, len(machine_torque_curve) + 1)]
# machine_torque["主轴扭矩曲线"] = [point, machine_torque_curve]
# # print(machine_torque)
#

# machine_health = {"机床体检评分": []}
# machine_health_curve = machine_specifications["机床体检评分"].tolist()[1:]
# point = [i for i in range(1, len(machine_health_curve) + 1)]
# machine_health["机床体检评分"] = [point, machine_health_curve]
# print(machine_health)


### 默认图形库
# pie_list = [
#     [str(value) for value in row]
#     for row in pd.read_csv(os.path.join('conf', 'mock_pie_data.csv')).values.tolist()
# ]


# finish_list = [
#     ["43", "32", "21"],
#     ["646", "544", "124"],
#     ["44", "1", "0"],
#     ["65", "55", "3"],
#     ["21", "3", "1"],
# ]


# finish_rate_list = [
#     [str(value) for value in row]
#     for row in pd.read_csv(os.path.join('conf', 'mock_finish_rate_list.csv')).values.tolist()
# ]


# list_seven = [
#     [str(value) for value in row]
#     for row in pd.read_csv(os.path.join('conf', 'mock_list_seven.csv')).values.tolist()
# ]

# product_list = ["机匣", "细长轴", "轴承外圈"]

# product_list: list[str] = ["连杆", "连杆", "连杆"]

# product_id = ["WX15736372", "STE78382782", "OIU-9032390"]

# product_id_list = ["091123", "324242", "324423"]

# product_id_list = [  # TODO: 想要获取"091123"，获取到"91123"，会自动将前导0去除
#     [str(value) for value in row]
#     for row in pd.read_csv(os.path.join('conf', 'mock_product.csv')).values.tolist()
# ]


# carbon_emissions = [
#     {
#         "day-4hour": ["200.0", "100.0", "200.0", "343.0", "250.0", "112.0"],
#         "week-1day": ["200.0", "100.0", "200.0", "343.0", "250.0", "112.0", "100.0"],
#         "month-5day": ["200.0", "100.0", "200.0", "343.0", "250.0", "112.0"],
#     },
#     {
#         "day-4hour": ["125.3", "563.3", "85.3", "343.0", "250.0", "112.0"],
#         "week-1day": ["200.0", "223.6", "365.3", "343.0", "250.0", "112.0", "100.0"],
#         "month-5day": ["200.0", "36.9", "200.0", "343.0", "169.3", "112.0"],
#     },
# ]
# carbon_per_money = [
#     {
#         "day-4hour": ["200.0", "100.0", "200.0", "343.0", "250.0", "112.0"],
#         "week-1day": ["200.0", "100.0", "200.0", "343.0", "250.0", "112.0", "100.0"],
#         "month-5day": ["200.0", "100.0", "200.0", "343.0", "250.0", "112.0"],
#     },
#     {
#         "day-4hour": ["200.0", "100.0", "200.0", "343.0", "250.0", "112.0"],
#         "week-1day": ["589.3", "100.0", "36.9", "343.0", "59.4", "112.0", "100.0"],
#         "month-5day": ["200.0", "639.3", "200.0", "269.3", "250.0", "112.0"],
#     },
# ]
# humidity = [
#     {
#         "day-4hour": ["0.98", "0.25", "0.68", "0.67", "0.97", "0.39"],
#         "week-1day": ["0.36", "0.37", "0.47", "0.59", "0.37", "0.76", "0.71"],
#         "month-5day": ["0.75", "0.39", "0.73", "0.37", "0.74", "0.98"],
#     },
#     {
#         "day-4hour": ["0.14", "0.678", "0.68", "0.67", "0.74", "0.39"],
#         "week-1day": ["0.36", "0.07", "0.47", "0.59", "0.59", "0.76", "0.04"],
#         "month-5day": ["0.75", "0.39", "0.78", "0.96", "0.74", "0.98"],
#     },
# ]
# temperature = [
#     {
#         "day-4hour": ["23.9", "32.7", "12.6", "14.7", "22.9", "10.8"],
#         "week-1day": ["23.9", "16.9", "29.0", "33.0", "20.0", "11.0", "10.0"],
#         "month-5day": ["20.0", "10.0", "20.0", "33.0", "20.0", "12.0"],
#     },
#     {
#         "day-4hour": ["43.9", "32.7", "32.6", "14.7", "22.9", "10.8"],
#         "week-1day": ["23.9", "16.9", "19.0", "33.0", "23.0", "11.0", "16.0"],
#         "month-5day": ["20.0", "10.0", "29.0", "38.0", "21.0", "12.0"],
#     },
# ]
# Electricity_consumption = [
#     {
#         "day-4hour": [["243.0", "123.0", "780.0", "398.0", "240.0", "113.0"],
#                       ["670.0", "300.6", "500.0", "300.0", "230.0", "145.0"]],
#         "week-1day": [["243.0", "123.0", "780.0", "398.0", "240.0", "113.0", "124.9"],
#                       ["670.0", "300.0", "500.0", "300.0", "230.0", "145.0", "124.9"]],
#         "month-5day": [["243.0", "123.0", "780.0", "398.0", "240.0", "113.0"],
#                        ["670.0", "300.0", "500.0", "300.0", "230.0", "145.0"]],
#     },
#     {
#         "day-4hour": [["234.0", "123.0", "780.0", "278.0", "240.0", "563.0"],
#                       ["567.0", "300.6", "345.0", "279.0", "230.0", "89.0"]],
#         "week-1day": [["243.0", "123.0", "780.0", "398.63", "240.0", "113.67", "124.9"],
#                       ["346.0", "300.0", "500.0", "764.0", "230.0", "145.0", "124.9"]],
#         "month-5day": [["243.0", "123.0", "349.0", "398.0", "167.0", "196.0"],
#                        ["670.0", "300.0", "279.0", "300.0", "368.0", "198.0"]],
#     },
# ]


# df_carbon_emissions = pd.read_csv(os.path.join('conf', 'mock_carbon_per_money.csv'))
# carbon_emissions = []
# chunks = [df_carbon_emissions.iloc[0:6]]
#
# # 将DataFrame转换为字典
# for chunk in chunks:
#     row_dict = {}
#     for column in chunk.columns:
#         row_dict[column] = chunk[column].tolist()
#     carbon_emissions.append(row_dict)
# # 将carbon_emissions列表中的数值转换为字符串
# for entry in carbon_emissions:
#     for key, value in entry.items():
#         entry[key] = [str(i) for i in value]
# # print("carbon_emissions:", carbon_emissions)
#


# df_carbon_per_money = pd.read_csv(os.path.join('conf', 'mock_carbon_emissions.csv'))
# carbon_per_money = []
# # 将CSV数据分为两部分，每6行一组
# chunks = [df_carbon_per_money.iloc[0:6]]
#
# # 将DataFrame转换为字典
# for chunk in chunks:
#     row_dict = {}
#     for column in chunk.columns:
#         row_dict[column] = chunk[column].tolist()
#     carbon_per_money.append(row_dict)
# for entry in carbon_per_money:
#     for key, value in entry.items():
#         entry[key] = [str(i) for i in value]
# # print("carbon_per_money: ", carbon_per_money)


# df_rank = pd.read_csv(os.path.join('conf', 'mock_nergy_rank.csv'))
# rank = []
# for _, row in df_rank.iterrows():
#     rank.append([list(df_rank.columns), list(row)])
# # print(rank)


# df_energy_partion = pd.read_csv(os.path.join('conf', 'mock_energy_partion.csv'))
# energy_partion = []
# # 将DataFrame转换为字典并添加到列表中
# for _, row in df_energy_partion.iterrows():
#     energy_partion.append(dict(zip(df_energy_partion.columns, row)))


# print(energy_partion)


# df_humidity = pd.read_csv(os.path.join('conf', 'mock_humidity.csv'), dtype=str)
# humidity = []
# # 将CSV数据分为两部分，每6行一组
# chunks = [df_humidity.iloc[0:6]]
#
# # 将DataFrame转换为字典
# for chunk in chunks:
#     row_dict = {}
#     for column in chunk.columns:
#         row_dict[column] = chunk[column].tolist()
#     humidity.append(row_dict)
# # print("humidity: ", humidity)


# df_temperature = pd.read_csv(os.path.join('conf', 'mock_temperature.csv'), dtype=str)
# temperature = []
# # 将CSV数据分为两部分，每6行一组
# chunks = [df_temperature.iloc[0:6]]
#
# # 将DataFrame转换为字典
# for chunk in chunks:
#     row_dict = {}
#     for column in chunk.columns:
#         row_dict[column] = chunk[column].tolist()
#     temperature.append(row_dict)
# # print("temperature:", temperature)


# # 读取CSV文件
# df = pd.read_csv(os.path.join('conf', 'mock_Electricity_consumption.csv'), dtype=str)
#
# # 初始化空列表以存储重构的数据
# Electricity_consumption = []
#
# # 按照原始数据结构，每12行分为一个数据条目
# # for i in range(0, len(df), 12):
# entry = {}
#
# # 处理day-4hour
# entry['day-4hour'] = [df.iloc[0:6]['day-4hour'].values.tolist(),
#                       df.iloc[6:12]['day-4hour'].values.tolist()]
#
# # 处理week-1day
# entry['week-1day'] = [df.iloc[0:6]['week-1day'].values.tolist(),
#                       df.iloc[6:12]['week-1day'].values.tolist()]
#
# # 处理month-5day
# entry['month-5day'] = [df.iloc[0:6]['month-5day'].values.tolist(),
#                        df.iloc[6:12]['month-5day'].values.tolist()]
#
# Electricity_consumption.append(entry)
#
# # 输出重构后的数据结构
# # print("Electricity_consumption: ", Electricity_consumption)


# df_product = pd.read_csv(os.path.join('conf', 'mock_product.csv'))

# 假设df是已经加载了CSV数据的DataFrame
# product = df_product.loc[0].to_dict()
# product_id = df_product.loc[1].to_dict()
# product_good_rate = df_product.loc[2].to_dict()

# 确保列名正确对应到目标变量的键上
# products = {col: product[col] for col in df_product.columns}
# product_ids = {col: product_id[col] for col in df_product.columns}
# product_good_rates = {col: product_good_rate[col] for col in df_product.columns}

# print(product)
# print(product_id)
# print(product_good_rate)


# 静态模拟百分比数据
# df_completion = pd.read_csv(os.path.join('conf', 'mock_completion_rate.csv'))
# run_status = df_completion.loc[0].to_dict()
# produce_finish_status = df_completion.loc[1].to_dict()
# # print(run_status)
# # print(produce_finish_status)


def __check_time_interval(time_interval: str) -> bool:
    """验证时间格式"""

    if time_interval == "day-4hour" or time_interval == "week-1day" or time_interval == "month-5day" or time_interval == "season-30d":
        return True
    else:
        return False


def __get_machine_num() -> int:
    """获取机床数量"""

    # MACHINES_CONFIG = os.path.join(str(path_depend), MACHINE_CONFIG)
    # machines_config = get_conf_file(MACHINE_CONFIG)
    # domTree = parse(machines_config)
    # rootNode = domTree.documentElement
    # machine = rootNode.getElementsByTagName("Machine")
    machine_num = len(machine)

    return machine_num


def __get_machine_id_list() -> list:
    """读取所有机床的id"""

    # MACHINES_CONFIG = os.path.join(str(path_depend), MACHINE_CONFIG)
    # machines_config = get_conf_file(MACHINE_CONFIG)
    # domTree = parse(machines_config)
    # rootNode = domTree.documentElement
    # machine = rootNode.getElementsByTagName("Machine")
    # machine_num = len(machine)

    machine_id_list = []
    for i in machine:
        machine_id = i.getElementsByTagName("id")[0].childNodes[0].data
        machine_id_list.append(machine_id)

    return machine_id_list


def __get_machine_name_list() -> list:
    """读取所有机床的名称"""

    # MACHINES_CONFIG = os.path.join(str(path_depend), MACHINE_CONFIG)
    # machine_num = len(machine)

    machine_name_list = []
    for i in machine:
        static = i.getElementsByTagName("Static")[0]
        machine_name = static.getElementsByTagName("machine_name")[0].childNodes[0].data
        machine_name_list.append(machine_name)

    return machine_name_list


def __gen_datetime_interval(time_interval: str):
    """计算时间区间"""

    end_time = datetime.now()

    start_time: datetime

    if time_interval == "day-4hour":
        start_time = end_time - timedelta(days=1)
    elif time_interval == "week-1day":
        start_time = end_time - timedelta(days=7)
    elif time_interval == "month-5day":
        start_time = end_time - timedelta(days=30)
    # elif time_interval == "season-30d":
    else:
        start_time = end_time - timedelta(days=90)

    start_timestamp = time.mktime(start_time.timetuple())
    start_time = datetime.fromtimestamp(floor(start_timestamp / 3600) * 3600)

    return (start_time, end_time)


def __select_log_from_machine(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    """查询每台机床的日志"""

    rows = select_from_single_machine(
        machine_index=machine_id,
        table_name="log",
        column_name=["id"],
        where=f"time>='{start_time}' and time<='{end_time}'",
    )
    # print("rows: ", rows)

    columns = ["id", "timestamp", "status"]

    log_data = pd.DataFrame()

    # 获取第一条数据
    # 处理数据为空的情况
    if not rows:
        # 符合时间要求的数据为空自动生成两行start_time,end_time
        log_data = pd.DataFrame(np.empty((2, 3)), columns=["id", "timestamp", "status"])

        # new_row = pd.Series([str(0), start_time.strftime("%Y-%m-%d %H:%M:%S"), MACHINE_STATUS["breakdown"]])

        # log_data.iloc[0, :] = [
        #     str(0),
        #     start_time.strftime("%Y-%m-%d %H:%M:%S"),
        #     MACHINE_STATUS["breakdown"],
        # ]
        # log_data.iloc[1, :] = [
        #     str(1),
        #     end_time.strftime("%Y-%m-%d %H:%M:%S"),
        #     MACHINE_STATUS["breakdown"],
        # ]

        log_data = pd.DataFrame(
            data=[
                [str(0), date_format(start_time), MACHINE_STATUS["breakdown"]],
                [str(1), date_format(end_time), MACHINE_STATUS["breakdown"]],
            ],
            columns=columns,
        )

    else:
        first_row_id = rows[0][0]
        if isinstance(first_row_id, (int, str, float)):
            first_row_id = int(first_row_id)
            if first_row_id != 1:
                # 符合时间要求的数据不为空, 且第一行不是数据库的第一行,
                # 则向前多取一行, 时间改为start_time, 最后一行改为 NOT_CONNECT,
                # 再插入时间为end_time,状态为NOT_CONNECT
                prev_rows = select_from_single_machine(
                    machine_index=machine_id,
                    table_name="log",
                    column_name=["id", "time", "status"],
                    where=f"id>='{first_row_id - 1}' and id<='{rows[len(rows) - 1][0]}'",
                )
                # 数据库的id可能不是从1开始的
                if prev_rows:
                    log_data = pd.DataFrame(prev_rows, columns=columns)
                    # 修改第一行的开始时间格式
                    log_data.iloc[0, 1] = date_format(start_time)
                    # 修改最后一行的状态
                    log_data.iloc[-1, 2] = MACHINE_STATUS["breakdown"]
                    # 待插入数据, 手动生成的最后一行
                    # d1 = pd.DataFrame(np.empty((1, 3)), columns=columns)
                    # d1.loc[0, :] = [0, date_format(end_time), MACHINE_STATUS["breakdown"]]
                    # 6-19 修改future warning错误
                    d1 = pd.DataFrame([[0, date_format(end_time), MACHINE_STATUS["breakdown"]]],
                                      columns=columns)

                    log_data = pd.concat([log_data, d1]).copy()

            else:
                # 符合时间要求的数据不为空, 但第一行为数据库的第一行, 取出数据后在第一行位置插入一行,
                # 时间改为start_time, 状态为NOT_CONNECT, 最后一行改为 NOT_CONNECT,
                # 再插入时间为end_time,状态为NOT_CONNECT
                log_data = pd.DataFrame(
                    select_from_single_machine(
                        machine_index=machine_id,
                        table_name="log",
                        column_name=["id", "time", "status"],
                        where=f"id>='{first_row_id}' and id<='{rows[len(rows) - 1][0]}'",
                    ),
                )
                log_data.columns = ["id", "timestamp", "status"]
                # 待插入数据
                d1 = pd.DataFrame(np.empty((1, 3)), columns=["id", "timestamp", "status"])
                d1.loc[0, :] = [0, date_format(start_time), MACHINE_STATUS["breakdown"]]
                log_data = pd.concat([d1, log_data]).copy()
                log_data.iloc[-1, 2] = MACHINE_STATUS["breakdown"]
                # 待插入数据
                d1 = pd.DataFrame(np.empty((1, 3)), columns=["id", "timestamp", "status"])
                d1.loc[0, :] = [0, date_format(end_time), MACHINE_STATUS["breakdown"]]
                log_data = pd.concat([log_data, d1]).copy()

    log_data["timestamp"] = log_data["timestamp"].apply(datetime.fromisoformat)
    return log_data


def __select_status_from_machine(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    status = select_from_single_machine(
        machine_index=machine_id,
        table_name="log ORDER BY id DESC LIMIT 1",
        column_name=["status"],
    )
    if not status:
        return None
    return str(status[0])


def __select_top_log_from_machine(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    where_condition = "start_time>='%s' and start_time<='%s'" % (start_time, end_time)
    log_id = select_from_single_machine(machine_index=machine_id, table_name="program", column_name=["id"],
                                        where=where_condition)
    if not log_id:
        return None
    if log_id == []:
        # 符合要求的数据为空
        return str(0), None
    else:
        # 存在时间区间内的数据
        where_condition = "id>='%s' and id<='%s'" % (log_id[0][0], log_id[len(log_id) - 1][0])
        log_data = pd.DataFrame(select_from_single_machine(machine_index=machine_id, table_name="program",
                                                           column_name=["name", "start_time", "end_time", "run_time"],
                                                           where=where_condition))
        log_data.columns = ["program_name", "start_time", "end_time", "run_time"]
        # log_data['start_time'] = log_data['start_time'].apply(lambda x: datetime.fromisoformat(x))
        # log_data['end_time'] = log_data['end_time'].apply(lambda x: datetime.fromisoformat(x))
    return str(1), log_data


def __select_program_count_from_machine(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    """
    根据log表 排除 运行程序名 相邻且连续的记录数 后剩下的记录数作为某台机床的完成零件数
    """
    where_condition = "time>='%s' and time<='%s'" % (start_time, end_time)
    data = select_from_single_machine(machine_index=machine_id, table_name="log", column_name=["id", "name"],
                                      where=where_condition)
    df = pd.DataFrame(data, columns=["id", "name"])
    # print("log数据：", data)
    # print(f"{machine_id}机床的 df 数据：", df)
    # 检查当前行的name是否与前一行不同
    df['is_unique'] = df['name'] != df['name'].shift()

    # 计算非连续重复行的数量
    program_count = len(df[df['is_unique']])
    return program_count


def __select_pro_from_machine(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    """根据开始结束时间筛选符合条件的program表id,name列"""
    # TODO：6-17 修改判断条件（run_time是否需要大于50？）
    where_condition_3 = f"end_time IS NOT NULL AND run_time > 50 AND start_time>='{start_time}' and start_time<='{end_time}'"
    log_id = select_from_single_machine(
        machine_index=machine_id,
        table_name="program",
        column_name=["id"],
        where=where_condition_3,
    )
    if log_id is None:
        return None

    if log_id == []:
        # program_data = pd.DataFrame(np.empty((1, 2)), columns=["id", "name"])
        # program_data.iloc[0, :] = [1, "nothing"]

        program_data = pd.DataFrame([[1, "nothing"]], columns=["id", "name"])

    else:
        where_condition_4 = "id>='%s' and id<='%s'" % (
            log_id[0][0], log_id[len(log_id) - 1][0])
        program_data = pd.DataFrame(
            select_from_single_machine(machine_index=machine_id, table_name="program", column_name=["id", "name"],
                                       where=where_condition_4))
    program_data.columns = ["id", "name"]

    return program_data


def __select_tool_from_machine(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    """从tool表中查找需要的数据"""

    # 根据开始结束时间筛选符合条件的program表id,name列
    # rows = select_from_single_machine(
    #     machine_index=machine_id,
    #     table_name="tool",
    #     column_name=["id", "number"],
    #     where=f"tool_start_time>='{start_time}' and tool_start_time<='{end_time}'",
    # )

    sql = """
        SELECT `tool`.`id`, IFNULL(`tm`.`tool_id`, `tool`.`tool_id`)
        FROM `tool`
        LEFT JOIN `tool_mapping` AS `tm` ON `tm`.`tool_no` = `tool`.`tool_id`
        WHERE `tool`.`tool_start_time` >= %(start_time)s AND `tool`.`tool_start_time` <= %(end_time)s
        AND `tool`.`tool_id` != 0 
    """

    rows = select_single_machine_with_raw(
        machine_index=machine_id,
        raw=sql,
        params={
            "start_time": start_time,
            "end_time": end_time,
        },
    )

    if not rows:
        # 没有数据, 或者查询失败
        return pd.DataFrame([[1, "nothing"]], columns=["id", "name"])

    return pd.DataFrame(rows, columns=["id", "name"])


def __log_count_runtime(log_runtime: pd.DataFrame) -> dict:
    # 统计log中四种状态RUNNING,IDLE,PAUSING,NOT_CONNECT所占时长和百分比
    # 统计runtime
    log_runtime["runtime"] = log_runtime.timestamp.diff()
    log_runtime["runtime_1"] = log_runtime["runtime"].shift(-1)
    log_runtime_1 = log_runtime.drop(columns="runtime")
    log_runtime_2 = log_runtime_1.dropna()
    df = log_runtime_2.groupby(["status"])["runtime_1"].sum()

    status = list(df.index)
    runtime = list(df.values)

    all_time = log_runtime.iloc[-1, 1] - log_runtime.iloc[0, 1]
    zero = pd.to_timedelta(0)

    if MACHINE_STATUS["running"] not in status:
        status.insert(0, MACHINE_STATUS["running"])
        runtime.insert(0, zero)

    if MACHINE_STATUS["shutdown"] not in status:
        status.insert(1, MACHINE_STATUS["shutdown"])
        runtime.insert(1, zero)

    if MACHINE_STATUS["pausing"] not in status:
        status.insert(2, MACHINE_STATUS["pausing"])
        runtime.insert(2, zero)

    if MACHINE_STATUS["breakdown"] not in status:
        status.insert(3, MACHINE_STATUS["breakdown"])
        runtime.insert(3, all_time - runtime[0] - runtime[1] - runtime[2])
    else:
        runtime[3] = all_time - runtime[0] - runtime[1] - runtime[2]

    runtime = [*map(lambda x: x / np.timedelta64(1, "h"), runtime)]
    all_time_1 = all_time / np.timedelta64(1, "h")

    # 计算百分比
    per = []
    for _, v in enumerate(runtime):
        per.append((v / all_time_1) * 100)

    d = {
        status[0]: [float("%.1f" % (runtime[0])), ("%.1f" % per[0]) + "%"],
        status[1]: [float("%.1f" % (runtime[1])), ("%.1f" % per[1]) + "%"],
        status[2]: [float("%.1f" % (runtime[2])), ("%.1f" % per[2]) + "%"],
        status[3]: [abs(float("%.1f" % (all_time_1 - runtime[0] - runtime[1] - runtime[2]))),
                    ("%.1f" % abs(100 - per[0] - per[1] - per[2])) + "%"],
    }

    return d


def __log_count_slice(log_df: pd.DataFrame, standar_range: str) -> dict:
    # 统计log切片中三种状态RUNNING,IDLE,PAUSING所占时长
    # 统计runtime
    s_range = standar_range
    log_runtime = log_df
    log_runtime["runtime"] = log_runtime.timestamp.diff()
    log_runtime["runtime_1"] = log_runtime["runtime"].shift(-1)
    log_runtime_1 = log_runtime.drop(columns="runtime")
    log_runtime_2 = log_runtime_1.dropna()
    df = log_runtime_2.groupby(["status"])["runtime_1"].sum()
    status = list(df.index)
    runtime = list(df.values)
    all_time = log_runtime.iloc[-1, 1] - log_runtime.iloc[0, 1]
    zero = pd.to_timedelta(0)

    if MACHINE_STATUS["running"] not in status:
        status.insert(0, MACHINE_STATUS["running"])
        runtime.insert(0, zero)

    if MACHINE_STATUS["shutdown"] not in status:
        status.insert(1, MACHINE_STATUS["shutdown"])
        runtime.insert(1, zero)

    if MACHINE_STATUS["pausing"] not in status:
        status.insert(2, MACHINE_STATUS["pausing"])
        runtime.insert(2, zero)

    if MACHINE_STATUS["breakdown"] not in status:
        status.insert(3, MACHINE_STATUS["breakdown"])
        runtime.insert(3, all_time - runtime[0] - runtime[1] - runtime[2])
    else:
        runtime[3] = all_time - runtime[0] - runtime[1] - runtime[2]

    runtime = [*map(lambda x: x / np.timedelta64(1, s_range), runtime)]
    all_time_new = all_time / np.timedelta64(1, s_range)
    runtime_new = []
    runtime_new.append(int(runtime[0]))
    runtime_new.append(int(runtime[1]))
    runtime_new.append(int(runtime[2]))
    runtime_new.append(int(all_time_new - int(runtime[0]) - int(runtime[1]) - int(runtime[2])))

    # 组装
    d = {
        status[0]: runtime_new[0],
        status[1]: runtime_new[1],
        status[2]: runtime_new[2],
        status[3]: runtime_new[3],
    }
    return d


def __log_aggre_slice(log_df: pd.DataFrame, standar_range: str) -> list:
    # 统计log切片中三种状态RUNNING,IDLE,PAUSING所占时长
    # 统计runtime
    s_range = standar_range
    log_runtime = log_df
    log_runtime["runtime"] = log_runtime.timestamp.diff()
    log_runtime["runtime_1"] = log_runtime["runtime"].shift(-1)
    log_runtime_1 = log_runtime.drop(columns="runtime")
    log_runtime_2 = log_runtime_1.dropna()
    log_runtime_3 = log_runtime_2.drop(
        index=log_runtime_2[log_runtime_2["runtime_1"] == pd.to_timedelta("0 days 00:00:00")].index)

    status = list(log_runtime_3["status"])
    runtime = list(log_runtime_3["runtime_1"])
    runtime = [*map(lambda x: x / np.timedelta64(1, s_range), runtime)]

    slice_list = []
    for i, v in enumerate(status):
        slice_list.append([v, runtime[i]])
    return slice_list


def __pro_count_all(pro_df: pd.DataFrame):
    """统计program中程序总数量和运行次数前五名的程序名"""

    # 统计程序总数量
    # 检查DataFrame的第一行是否为空
    if pro_df.iloc[0, 1] == "nothing":
        pro_number = 0
        pro_names = [""]
        runtime = [0]
    else:
        # 计算程序总数量
        pro_number = pro_df.shape[0]
        # 根据程序名统计运行次数
        df = pro_df["name"].value_counts()

        # 获取运行次数前五的程序名和对应的次数，限制结果数量不超过5个
        len_pro = min(5, len(df))
        pro_names: list[str | int] = list(df.index)[:len_pro]

        # 截取程序名后半段
        # for i in range(len(pro_names)):
        # pro_names[i] = pro_names[i].split("/_N_",-1)[-1]
        pro_names = list(map(reduce_program_name, pro_names))

        # 将运行次数从浮点数转换为整数
        runtime = list(df.values)[:len_pro]
        for i, v in enumerate(runtime):
            runtime[i] = int(v)

    # 组装成dict，键是程序名，值是对应的运行次数
    nrs = zip(pro_names, runtime)
    # 记录了运行次数前五的程序名及其运行次数
    names_times = dict((pro_name, run_time) for pro_name, run_time in nrs)
    # 返回一个包含两个元素的元组：第一个元素是程序的总数量，第二个元素是names_times字典，记录了运行次数前五的程序名及其运行次数
    return pro_number, names_times


def __pro_count_number(pro_df: pd.DataFrame) -> int:
    # 统计program中程序总数量作为本日产量

    pro_data = pro_df
    # 统计程序总数量
    if pro_data.iloc[0, 1] == "nothing":
        pro_number = 0
        pro_names = [""]
        runtimes = [0]
    else:
        pro_number = pro_data.shape[0]

    return pro_number


def __gen_slice(log_data: pd.DataFrame) -> tuple:
    """
    1、 1h间隔：12-20hours
    2、 4h间隔：20-72hours
    3、 12h间隔：3-6days
    4、 1d间隔：6-25days
    5、 5d间隔：26-60days
    6、 30d间隔：61-900days
    """
    # 计算log_data中时间跨度
    period = log_data.iloc[-1, 1] - log_data.iloc[0, 1]
    standard = []
    # 定义了6个不同尺度的时间间隔（1小时、4小时、12小时、1天、5天、30天）
    standard.append(pd.to_timedelta("0 days 1:00:00"))
    standard.append(pd.to_timedelta("0 days 4:00:00"))
    standard.append(pd.to_timedelta("0 days 12:00:00"))
    standard.append(pd.to_timedelta("1 days 00:00:00"))
    standard.append(pd.to_timedelta("5 days 00:00:00"))
    standard.append(pd.to_timedelta("30 days 00:00:00"))
    period_max = []
    # 定义了不同尺度的时间间隔的周期上限
    period_max.append(pd.to_timedelta("0 days 20:00:00"))
    period_max.append(pd.to_timedelta("3 days 00:00:00"))
    period_max.append(pd.to_timedelta("6 days 00:00:00"))
    period_max.append(pd.to_timedelta("25 days 00:00:00"))
    period_max.append(pd.to_timedelta("60 days 00:00:00"))
    period_max.append(pd.to_timedelta("900 days 00:00:00"))
    # 映射关系r，关联每个间隔与其对应的标签
    r = {standard[0]: "m", standard[1]: "m", standard[2]: "m", standard[3]: "h", standard[4]: "h", standard[5]: "h"}
    # interval记录了 可以覆盖log_data的时间跨度的间隔
    interval = pd.to_timedelta("30 days 00:00:00")
    # 遍历预先定义的时间间隔上限，找到第一个能覆盖整个数据时间跨度的间隔，将其作为最终的时间间隔
    for i, v in enumerate(period_max):
        if period <= v:
            interval = standard[i]
            break
    standard_range = r[interval]

    dfs = []
    # 根据确定的时间间隔，计算需要分割出多少个时间片
    num = floor(period / interval)
    for i in range(num):
        pro_time = log_data.iloc[0, 1] + interval
        for j in range(log_data.shape[0] - 1):
            if pro_time >= log_data.iloc[j, 1] and pro_time <= log_data.iloc[j + 1, 1]:
                # 待插入行索引
                row_n = j + 1
                # 待插入数据
                # d1 = pd.DataFrame(np.empty((1, 3)), columns=["id", "timestamp", "status"])
                # d1.iloc[0, :] = [int(data_OEE.iloc[j, 0]) + 1, pro_time, data_OEE.iloc[j, 2]]
                # d2 = pd.DataFrame(np.empty((1, 3)), columns=["id", "timestamp", "status"])
                # d2.iloc[0, :] = [int(data_OEE.iloc[j + 1, 0]) - 1, pro_time, data_OEE.iloc[j, 2]]
                # 6-19 修改future warning错误
                d1 = pd.DataFrame([[int(log_data.iloc[j, 0]) + 1, pro_time, log_data.iloc[j, 2]]],
                                  columns=["id", "timestamp", "status"])

                d2 = pd.DataFrame([[int(log_data.iloc[j + 1, 0]) - 1, pro_time, log_data.iloc[j, 2]]],
                                  columns=["id", "timestamp", "status"])

                # 拆分
                pd_arr1 = log_data.iloc[:row_n]
                pd_arr2 = log_data.iloc[row_n:]
                # 参数：添加数据, 是否无视行索引
                pd_arr = pd.concat([pd_arr1, d1]).copy()
                log_data = pd.concat([d2, pd_arr2]).copy()
                dfs.append(pd_arr)
                break
    return (dfs, interval, standard_range)


def __gen_period(
        start_time: datetime,
        end_time: datetime,
        inter: pd.Timedelta,
) -> list:
    """根据开始时间, 结束时间, 时间间隔生成横坐标"""

    zero = pd.to_timedelta(0)
    period = end_time - start_time
    num = int(period / inter)
    if (period % inter) == zero:
        num = num - 1

    p = []
    t = start_time
    for i in range(num + 1):
        s = datetime.strftime(t, "%Y-%m-%d %H") + "时"
        t = t + inter
        p.append(s)

    return p


def __gen_interval(length: int) -> int:
    len = length
    if len <= 6:
        l = 0
    else:
        if len <= 12:
            l = 2
        else:
            l = 5
    return l


def _update_machine_OEE(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    # 1. 在数据库中筛选符合条件的 log 表记录
    data_log = __select_log_from_machine(
        machine_id=machine_id,
        start_time=start_time,
        end_time=end_time,
    )

    # print("data_log: ", data_log)

    if not isinstance(data_log, pd.DataFrame):
        return None
    # 1.2 runtime, 统计所有时间的四种状态的时长和百分比
    data_log_runtime = data_log.copy()
    runtime = __log_count_runtime(data_log_runtime)

    return runtime


def __update_machine_slice(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    # 1. 筛选符合条件的 log 记录
    data_log = __select_log_from_machine(
        machine_id=machine_id,
        start_time=start_time,
        end_time=end_time,
    )
    if data_log is None:
        return None

    # 1.3 统计切片数据
    data_log_slice = data_log.copy()
    series = {"RUNNING": [], "IDLE": [], "PAUSING": [], "NOT_CONNECT": []}
    strip_list = __gen_slice(data_log_slice)
    log_list = strip_list[0]
    inter = strip_list[1]
    s_range = strip_list[2]

    # period = __gen_period(start_time,end_time,inter)
    # interval = __gen_interval(len(period))

    runtime_slice = []
    # 统计
    for i in range(len(log_list)):
        slice_list = __log_aggre_slice(log_list[i], s_range)
        for j in range(len(slice_list)):
            runtime_slice.append(slice_list[j])
    return runtime_slice


def _update_machine_slice_daily(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    """1. 筛选符合条件的 log 记录"""
    data_log = __select_log_from_machine(
        machine_id=machine_id,
        start_time=start_time,
        end_time=end_time,
    )
    if data_log is None:
        return None

    # 1.3 统计切片数据
    data_log_slice = data_log.copy()
    strip_list = __gen_slice(data_log_slice)
    log_list = strip_list[0]
    inter = strip_list[1]
    s_range = strip_list[2]

    # period = __gen_period(start_time,end_time,inter)
    # interval = __gen_interval(len(period))
    # 统计
    series = {"RUNNING": [], "IDLE": [], "PAUSING": [], "NOT_CONNECT": []}
    for i in range(len(log_list)):
        runtime_slic = __log_count_slice(log_list[i], s_range)
        series["RUNNING"].append(runtime_slic["RUNNING"])
        series["IDLE"].append(runtime_slic["IDLE"])
        series["PAUSING"].append(runtime_slic["PAUSING"])
        series["NOT_CONNECT"].append(runtime_slic["NOT_CONNECT"])
    return series


def _update_machine_prog_rank(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    # 2 program表统计：
    # 2.1 根据开始结束时间筛选符合条件的program表id,name列
    program_data = __select_pro_from_machine(machine_id, start_time, end_time)
    if not isinstance(program_data, pd.DataFrame):
        return None

    # 2.2 统计时间内程序总数量 和 运行次数前五名的字典{ 程序名:运行次数 }
    pro_number, run_number = __pro_count_all(program_data)
    # 返回两个元素：第一个元素是程序的总数量，第二个元素是names_times字典，记录了运行次数前五的程序名及其运行次数
    return pro_number, run_number


def _update_machine_tool_rank(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    """计算刀具使用排行"""
    ## 2 program表统计：
    # 2.1 根据开始结束时间筛选符合条件的tool表id,name列
    tool_data = __select_tool_from_machine(machine_id, start_time, end_time)

    # 2.2 统计时间内程序数量和运行次数前五名的程序名
    tool_times, tool_nummber = __pro_count_all(tool_data)

    return tool_nummber


def _update_machine_finish_number(
        machine_id: int,
        start_time: datetime,
        end_time: datetime,
):
    ###
    program_data = __select_pro_from_machine(machine_id, start_time, end_time)
    # print(f'{program_data}\n')
    if not isinstance(program_data, pd.DataFrame):
        return None
    # print(program_data)
    finish_number = __pro_count_number(program_data)

    return finish_number


def __divide_and_average(data, num_parts=6):
    """
    将传入的data数据分组，并计算每组的平均值
    :param data: 从数据库中查询到的排序好的元组列表，其中每个元组包含一整行数据
    :param num_parts: 分组的组数
    :return: 分组的组数个平均值数据组成的列表
    """
    # 将数据转换为浮点数数组
    data_float = np.array([float(item[0]) for item in data])

    # 将数据均匀分割成num_parts份
    split_data = np.array_split(data_float, num_parts)

    # 对每个子集计算平均值
    averages = []
    for subset in split_data:
        avg = np.mean(subset, axis=0)  # 计算平均值，axis=0表示按列计算
        averages.append(avg)

    return averages


def __select_temperature_from_machine(start_time, end_time):
    """
    从H8000机床获取温度数据
    """
    return select_from_single_machine(
        machine_index=801020,
        table_name="temperature",
        column_name=["temperature"],
        where=f"create_time >= '{start_time}' and create_time <= '{end_time}' order by create_time asc",
    )


def __select_humidity_from_machine(start_time, end_time):
    """
    从H8000机床获取湿度数据
    """
    return select_from_single_machine(
        machine_index=801020,
        table_name="humidity",
        column_name=["humidity"],
        where=f"create_time >= '{start_time}' and create_time <= '{end_time}' order by create_time asc",
    )


def __select_electricity_from_machine(start_time, end_time):
    """
    从H8000机床获取总电量数据
    """
    return select_from_single_machine(
        machine_index=801020,
        table_name="electricity",
        column_name=["electricity"],
        where=f"create_time >= '{start_time}' and create_time <= '{end_time}' order by create_time asc",
    )


def _get_temperature_avg_data_from_machine(start_time, end_time):
    """
    获取时间范围内的6组温度的平均值数据
    """
    try:
        data = __select_temperature_from_machine(start_time, end_time)

        res = __divide_and_average(data, num_parts=6)

        return [round(r, 2) for r in res]
    except Exception as e:
        print(f"获取温度数据失败：{e}")
        return []


def _get_humidity_avg_data_from_machine(start_time, end_time):
    """
    获取时间范围内的6组湿度的平均值数据
    """
    try:
        data = __select_humidity_from_machine(start_time, end_time)

        averages = __divide_and_average(data, num_parts=6)

        return [round(average, 2) for average in averages]
    except Exception as e:
        print(f"获取湿度数据失败：{e}")
        return []


def _get_electricity_avg_data_from_machine(start_time, end_time):
    """
    获取时间范围内的6组总电量的平均值数据，作为总排放量/t
    """
    try:
        data = __select_electricity_from_machine(start_time, end_time)

        averages = __divide_and_average(data, num_parts=6)

        res = [round(average * 0.000785, 2) for average in averages]

        # 防止数据变化太少，导致前端展示的曲线为一根直线
        res[0] = round(res[0] - res[0] * 0.1, 2)
        res[-1] = round(res[-1] + res[-1] * 0.1, 2)
        return res
    except Exception as e:
        print(f"获取总排放量数据失败：{e}")
        return []


def _get_electricity_diff_data_from_machine(start_time, end_time):
    """
    从H8000机床 获取时间范围内的7组总电量的平均值，返回每组平均值之间的差值，作为每日排放量/kg
    """
    try:
        data = __select_electricity_from_machine(start_time, end_time)

        averages = __divide_and_average(data, num_parts=7)

        diff_data = [round((b - a) * 0.785, 2) for a, b in zip(averages[:-1], averages[1:])]

        return diff_data
    except Exception as e:
        print(f"获取每日排放量数据失败：{e}")
        return []
    # return [round(dd*0.785, 2) for dd in diff_data]


def UpdAllOEE(time_interval: str):
    """2-1 更新所有机床的OEE    车间总览 - 设备有效运行率； 车间大屏 - 右上角的设备有效运行率"""

    temp_dict = {}
    dict_list = [0, 0, 0, 0]
    # TODO： 6-17 将shutdown的位置放到最后（原来：["running", "shutdown", "breakdown", "pausing"]）
    status_list = ["running", "breakdown", "pausing", "shutdown"]
    per_list = []

    if __check_time_interval(time_interval=time_interval):
        machine_num = __get_machine_num()
        machine_id_list = __get_machine_id_list()
        start_time, end_time = __gen_datetime_interval(time_interval=time_interval)

        # print(start_time, end_time)

        for i in range(machine_num):
            OEE_list = _update_machine_OEE(
                int(machine_id_list[i]),
                start_time,
                end_time,
            )
            # print("OEE_list: ", OEE_list)
            if OEE_list is None:
                return None

            for j in range(len(status_list)):
                dict_list[j] += OEE_list[MACHINE_STATUS[status_list[j]]][0]
                # TODO: 添加 保留小数点后2位 操作
        per_list = [str(round(x / sum(dict_list), 2)) for x in dict_list]

        # 7-1 防止per_list列表总和小于1
        per_list_int = [float(x) for x in per_list]
        if sum(per_list_int) < 1:
            difference = 1 - sum(per_list_int)
            # 如果总和少于1，则将缺少的部分填充到running状态中，使其的总和保证 = 1
            per_list[0] = str(float(per_list[0]) + difference)

        # 防止per_list列表总和大于1
        if sum(per_list_int) > 1:
            difference = sum(per_list_int) - 1
            # 如果总和大于1，则将多余的部分从pausing状态中删除，使其的总和保证 = 1
            per_list[2] = str(float(per_list[2]) - difference)

        # 7-2 防止出现有些状态的百分比大于100，有些百分比小于0
        for per in per_list:
            if float(per) > 1 or float(per) < 0:
                per_list = ["0.4", "0.0", "0.6", "0.0"]
                break

    temp_dict["task"] = str("update_all_OEE")
    temp_dict["time_interval"] = str(time_interval)
    if sum(dict_list) != 0:
        temp_dict["pie_list"] = [status_list, per_list]
        j = json.dumps(temp_dict, ensure_ascii=False)
    else:
        temp_dict["pie_list"] = "时间间隔不正确, 请保证时间间隔是day-4hour, week-1day, month-5day之一"
        j = json.dumps(temp_dict, ensure_ascii=False)
    # print("更新所有机床OEE：", j)
    publish_message(LOW_FREQ, j)

    return per_list


def UpdFinish(time_interval: str) -> None:
    """2-2 更新所有机床完成件数   车间总览 - 设备完成排名 - 设备完成数  【更新后此接口作废】"""
    temp_dict = {}
    prog_list = []
    machine_list = []
    if __check_time_interval(time_interval=time_interval):
        machine_num = __get_machine_num()
        machine_id_list = __get_machine_id_list()
        machine_name_list = __get_machine_name_list()
        for i in range(machine_num):
            start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
            run_number = _update_machine_finish_number(machine_id_list[i], start_time, end_time)
            # print(run_number)
            if run_number is None:
                return None
            machine_name = machine_name_list[i]
            prog_list.append(run_number)
            machine_list.append(machine_name)

    temp_dict["task"] = str("update_finish")
    temp_dict["time_interval"] = str(time_interval)
    temp_dict["finish_list"] = [machine_list, prog_list]
    # print([machine_list, prog_list])
    j = json.dumps(temp_dict, ensure_ascii=False)

    publish_message(LOW_FREQ, j)


def UpdFinishRate(time_interval: str) -> None:
    """2-3 更新所有机床完成率, 车间总览 - 设备完成排名 - 设备完成率  【更新后此接口作废】"""

    temp_dict = {}
    prog_list = []
    machine_list = []
    if __check_time_interval(time_interval=time_interval):
        machine_num = __get_machine_num()
        machine_id_list = __get_machine_id_list()
        machine_name_list = __get_machine_name_list()
        for i in range(machine_num):
            start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
            finish_rate = 0
            run_number = _update_machine_finish_number(machine_id_list[i], start_time, end_time)
            # TODO: 修改完成率逻辑(完成零件数 除 运行程序数或执行任务数)，待验证
            rows = select_from_single_machine(
                machine_index=machine_id_list[i],
                table_name="program",
                column_name=["COUNT(id)"],
                where=f"end_time IS NOT NULL AND run_time > 50 AND start_time >= '{start_time}' AND start_time <= '{end_time}'",
            )
            if rows:
                if rows[0][0] != 0:
                    finish_rate = round(int(rows[0][0]) / run_number, 2)
            machine_name = machine_name_list[i]
            prog_list.append(str(finish_rate))
            machine_list.append(machine_name)

    temp_dict["task"] = str("update_finish_rate")
    temp_dict["time_interval"] = str(time_interval)
    # temp_dict["finish_rate_list"] = [machine_list, prog_list]
    temp_dict["finish_rate_list"] = [machine_list, prog_list]
    # print("--------------------------设备完成率: ", prog_list)
    j = json.dumps(temp_dict, ensure_ascii=False)

    publish_message(LOW_FREQ, j)


def UpdWholeIndex(time_interval: str) -> None:
    """2-4 车间总览 - 设备运行状态   【去除合格率】 """
    # dice = random.randint(0,1)
    # equip_index = [{"day-4hour":[['running', 'shutdown', 'breakdown', 'pausing'],  ['10', '20', '5', '5']],
    #         "week-1day":[['running', 'shutdown', 'breakdown', 'pausing'],  ['100', '200', '50', '50']],
    #         "month-5day":[['running', 'shutdown', 'breakdown', 'pausing'],  ['1000', '2000', '500', '500']],
    #         },
    #         {"day-4hour":[['running', 'shutdown', 'breakdown', 'pausing'],  ['25', '16', '3', '2']],
    #         "week-1day":[['running', 'shutdown', 'breakdown', 'pausing'],  ['229', '169', '32', '26']],
    #         "month-5day":[['running', 'shutdown', 'breakdown', 'pausing'],  ['596', '369', '58', '11']],
    #         }]
    # product = [{
    #     "day-4hour":'10',
    #     "week-1day":'100',
    #     "month-5day":'1000',
    # },
    # {
    #     "day-4hour":'85',
    #     "week-1day":'129',
    #     "month-5day":'2369',
    # }]
    # product_rate = [{
    #     "day-4hour":'0.48',
    #     "week-1day":'0.98',
    #     "month-5day":'0.67',
    # },
    # {
    #     "day-4hour":'0.59',
    #     "week-1day":'0.23',
    #     "month-5day":'0.98',
    # }]
    temp_dict = {}
    temp_dict["task"] = str("update_whole_index")

    status_index = ["running", "shutdown", "breakdown", "pausing"]
    machine_status_num = [0, 0, 0, 0]
    product_finish_number = 0
    if __check_time_interval(time_interval=time_interval):
        prog_list = []
        machine_num = __get_machine_num()
        machine_id_list = __get_machine_id_list()
        for i in range(machine_num):
            start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
            # 设备运行状态为 查询每个机床的log表中最新行的status列数据
            machine_status = __select_status_from_machine(machine_id_list[i], start_time, end_time)
            machine_status = machine_status.replace("(", "")
            machine_status = machine_status.replace(")", "")
            machine_status = machine_status.replace("[", "")
            machine_status = machine_status.replace("]", "")
            machine_status = machine_status.replace("'", "")
            machine_status = machine_status.replace('"', "")
            machine_status = machine_status.replace(",", "")

            index = status_index.index(MACHINE_STATUS_INV[machine_status])
            machine_status_num[index] = machine_status_num[index] + 1

            # run_number = _update_machine_finish_number(machine_id_list[i], start_time, end_time)
            # if run_number is None:
            #     return None
            # prog_list.append(run_number)

            # 新增：计算所有机床的完成零件数
            rows = select_from_single_machine(
                machine_index=machine_id_list[i],
                table_name="program",
                column_name=["COUNT(id)"],
                where=f"end_time IS NOT NULL AND run_time > 50 AND start_time >= '{start_time}' AND start_time <= '{end_time}'",
            )
            if rows:
                product_finish_number += int(rows[0][0])

    temp_dict["time_interval"] = str(time_interval)
    temp_dict["equip_index"] = [status_index, machine_status_num]
    # # # print([status_index, machine_status_num])
    # TODO：6-17 修改本日产量 && 修改合格率为手动配置
    # TODO： update 修改 本日产量 的获取逻辑  并且 删除product_rate完成率
    temp_dict["product"] = product_finish_number
    # temp_dict["product"] = sum(prog_list)

    # 通过所有机床完成零件数 / 所有机床运行的程序数 = 总完成率
    # temp_dict["product_rate"] = round(int(product_finish_number) / sum(prog_list), 2)

    # 6-18 修改完成率获取方式
    # temp_dict["product_rate"] = run_status[time_interval]

    j = json.dumps(temp_dict, ensure_ascii=False)
    publish_message(LOW_FREQ, j)


def UpdEnv(time_interval: str) -> None:
    """2-5 车间总览-温湿度和电量"""
    temp_dict = {"task": str("update_env")}
    start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
    try:
        temp_dict["time_interval"] = str(time_interval)
        # temp_dict["carbon_emissions"] = carbon_emissions[dice][time_interval]
        # temp_dict["carbon_per_money"] = carbon_per_money[dice][time_interval]
        # temp_dict["humidity"] = humidity[dice][time_interval]
        # temp_dict["temperature"] = temperature[dice][time_interval]
        # temp_dict["Electricity consumption"] = Electricity_consumption[dice][time_interval]
        # temp_dict["carbon_emissions"] = ["200.0", "100.0", "200.0", "343.0", "250.0", "112.0"]
        # temp_dict["carbon_per_money"] = ["200.0", "100.0", "200.0", "343.0", "250.0", "112.0"]
        # TODO: 最新版 修改接口 待测试
        temp_dict["humidity"] = _get_humidity_avg_data_from_machine(start_time, end_time)
        temp_dict["temperature"] = _get_temperature_avg_data_from_machine(start_time, end_time)
        temp_dict["daily_carbon_emissions"] = _get_electricity_diff_data_from_machine(start_time,
                                                                                      end_time)  # 每日排放量
        temp_dict["total_carbon_emissions"] = _get_electricity_avg_data_from_machine(start_time,
                                                                                     end_time)  # 总排放量
        j = json.dumps(temp_dict, ensure_ascii=False)
    except KeyError:
        temp_dict["carbon_emissions"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
        j = json.dumps(temp_dict, ensure_ascii=False)
    publish_message(LOW_FREQ, j)


# def UpdProduct(time_interval: str) -> None:
#     """2-6 车间总览 - old(车间经济指标) 生产完成情况【假数据】 """
#
#     product_finish_number = 0
#     if __check_time_interval(time_interval=time_interval):
#         start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
#         # print(start_time, end_time)
#         machine_id_list = __get_machine_id_list()
#
#         for m_id in machine_id_list:
#             rows = select_from_single_machine(
#                 machine_index=m_id,
#                 table_name="program",
#                 column_name=["COUNT(id)"],
#                 where=f"end_time IS NOT NULL AND run_time > 50 AND start_time >= '{start_time}' AND start_time <= '{end_time}'",
#             )
#             if rows:
#                 product_finish_number += int(rows[0][0])
#
#     # product_finish_number = {
#     #     "day-4hour": "65",
#     #     "week-1day": "44",
#     #     "month-5day": "63",
#     # }
#     # product_good_rate = {
#     #     "day-4hour": "1",
#     #     "week-1day": "1",
#     #     "month-5day": "1",
#     # }
#
#     temp_dict = {}
#     temp_dict["task"] = str("update_product")
#     try:
#         temp_dict["time_interval"] = str(time_interval)
#         temp_dict["product"] = "待删"
#         temp_dict["product_id"] = "待删"
#         # temp_dict["product_finish_number"] = product_finish_number[time_interval]
#         temp_dict["product_finish_number"] = product_finish_number
#         # TODO： update 添加product_planned_finish_number 删除product  product_id  product_good_rate   并且修改 已完成数 和 计划完成数 的处理逻辑
#         temp_dict["product_good_rate"] = produce_finish_status[time_interval]
#         j = json.dumps(temp_dict, ensure_ascii=False)
#     except KeyError:
#         temp_dict["product"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
#         j = json.dumps(temp_dict, ensure_ascii=False)
#     # print(j)
#     publish_message(LOW_FREQ, j)
def UpdProduct(time_interval: str) -> None:
    """2-6 车间总览 - 生产完成进度情况   【通过mes接口获取 已完成数 和 计划完成数 的数据】 """

    product_finish_number = 0
    product_planned_finish_number = 0
    if __check_time_interval(time_interval=time_interval):
        start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
        machine_name_list = __get_machine_name_list()
        # machine_id_list = __get_machine_id_list()
        url = "http://192.168.25.140:8081/mes-web/zcplanmanagerController/zyfhGxrw11.do"

        for m_name in machine_name_list:
            if m_name == "NHM6300":
                # 查询字符串参数
                params = {
                    "manuCode": "H6300",
                    "startDate": start_time.strftime("%Y-%m-%d"),
                    "endDate": end_time.strftime("%Y-%m-%d")
                }
            else:
                params = {
                    "manuCode": m_name,
                    "startDate": start_time.strftime("%Y-%m-%d"),
                    "endDate": end_time.strftime("%Y-%m-%d")
                }
            print(f"请求的params： ", params)
            # 发送POST请求，注意这里没有请求体数据，如果有请用data或json参数添加
            response = requests.get(url, params=params)

            # 检查响应状态码
            if response.status_code == 200:
                res = response.json()
                if not res:
                    print(f"{m_name}在{params['startDate']} ~ {params['endDate']}时间段内没有生产记录")
                    continue
                tmp = res[0]['SHULIANG']
                # 累加所有机床的完成数与计划完成数
                product_finish_number += int(tmp.split("/")[0])
                product_planned_finish_number += int(tmp.split("/")[1])
            else:
                print(f"请求失败，状态码：{response.status_code}")

    temp_dict = {}
    temp_dict["task"] = str("update_product")
    try:
        # TODO： update 添加product_planned_finish_number 删除product  product_id  product_good_rate   并且修改 已完成数 和 计划完成数 的处理逻辑
        temp_dict["time_interval"] = str(time_interval)
        # temp_dict["product_finish_number"] = 33
        temp_dict["product_finish_number"] = product_finish_number
        temp_dict["product_planned_finish_number"] = product_planned_finish_number
        # temp_dict["product_planned_finish_number"] = 66
        j = json.dumps(temp_dict, ensure_ascii=False)
    except KeyError:
        temp_dict["product"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
        j = json.dumps(temp_dict, ensure_ascii=False)
    print("已完成数和计划完成数： ", product_finish_number, product_planned_finish_number)
    publish_message(LOW_FREQ, j)


def UpdProgRank(time_interval: str, machine_id: int):
    """2-10 进入机床 - 程序运行次数排行； 机床动态 - 程序运行次数"""

    temp_dict = {}
    # 程序名称列表
    prog_list = []
    # 程序运行次数列表
    run_time_list = []

    if __check_time_interval(time_interval=time_interval):
        start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
        # 返回一个包含两个元素的元组：第一个元素是程序的总数量，第二个元素是names_times字典，记录了运行次数前五的程序名及其运行次数
        machine_prog_rank = _update_machine_prog_rank(machine_id, start_time, end_time)

        if not machine_prog_rank:
            return None

        prog_number, run_number = machine_prog_rank
        # key：程序名    value：运行次数
        for key, value in run_number.items():
            prog_list.append(key)
            run_time_list.append(value)

    temp_dict["task"] = str("update_programm_rank")
    temp_dict["time_interval"] = str(time_interval)
    temp_dict["machine_id"] = str(machine_id)
    if sum(run_time_list) != 0:
        temp_dict["programm_rank"] = [prog_list, run_time_list]
        j = json.dumps(temp_dict, ensure_ascii=False)
    else:
        temp_dict["programm_rank"] = "该时间段内没有运行过程序!"
        j = json.dumps(temp_dict, ensure_ascii=False)
    # print("程序排名： ", j)
    publish_message(LOW_FREQ, j)

    return temp_dict


# NOTE 接口已调整, 等待前端界面调整后测试
def UpdToolRank(time_interval: str, machine_id: int):
    """2-11 进入机床 - 刀具使用次数排行； 机床动态 - 刀具使用次数"""

    temp_dict = {}
    run_time_list = []
    prog_list = []

    if __check_time_interval(time_interval=time_interval):
        start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
        run_number = _update_machine_tool_rank(machine_id, start_time, end_time)

        for key, value in run_number.items():
            prog_list.append(key)
            run_time_list.append(value)

    temp_dict["task"] = str("update_tool_rank")
    temp_dict["time_interval"] = str(time_interval)
    temp_dict["machine_id"] = str(machine_id)

    if sum(run_time_list) != 0:
        temp_dict["tool_rank"] = [prog_list, run_time_list]
        j = json.dumps(temp_dict, ensure_ascii=False)
    else:
        temp_dict["tool_rank"] = "该时间段内没有使用过刀具!"
        j = json.dumps(temp_dict, ensure_ascii=False)
    # print("刀具排名： ", j)
    publish_message(LOW_FREQ, j)

    return temp_dict


def UpdLog(time_interval: str):
    """2-15 更新车间日志   【更新后此接口作废】"""
    # TODO: log 字段返回的类型应该是列表还是字典？目前是列表

    machine_num = __get_machine_num()
    machine_id_list = __get_machine_id_list()
    start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
    log = None
    log_dict = {}
    for i in range(machine_num):
        flag, log_data = __select_top_log_from_machine(
            int(machine_id_list[i]),
            start_time,
            end_time,
        )
        if flag is None:
            return None
        if flag == str(1):
            # print("log_data", log_data)
            log = pd.concat([log, log_data]).copy()
    try:
        # 在转换为字典之前，将Timestamp转换为字符串
        str_log = log.apply(lambda col: col.dt.strftime('%Y-%m-%d %H:%M:%S') if col.dtype == 'datetime64[ns]' else col)
        log_dict = str_log.to_dict(orient="records")
    except Exception as e:
        print(f"Error during conversion: {e}")
        log_dict = []

    temp_dict = {}
    temp_dict["task"] = str("update_workshop_log")
    temp_dict["time_interval"] = time_interval
    # print("log_dict:", log_dict)
    # print("log", log)
    # TODO:待测试,将log_dict改成log_dict[0]
    temp_dict["log"] = log_dict[0]

    j = json.dumps(temp_dict, ensure_ascii=False)
    publish_message(LOW_FREQ, j)

    return temp_dict


def UpdSlice(time_interval: str, machine_id: int):
    """2-13 更新单台机床切片  车间大屏 - 运行状态切片"""

    machine_name_list = __get_machine_name_list()
    machine_id_list = __get_machine_id_list()
    start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
    runtime_slice = __update_machine_slice(machine_id, start_time, end_time)

    if runtime_slice is None:
        return None

    slice_final = []

    for slice in runtime_slice:
        slice_final.append([slice[1], MACHINE_STATUS_INV[str(slice[0])]])

    temp_dict = {}
    temp_dict["task"] = str("update_equip_slice")
    temp_dict["machine_id"] = str(machine_id)
    temp_dict["machine_name"] = str(machine_name_list[int(machine_id_list.index(str(machine_id)))])
    temp_dict["time_interval"] = str(time_interval)
    temp_dict["machine_slice"] = slice_final
    j = json.dumps(temp_dict, ensure_ascii=False)
    # print("运行状态切片: ", j)
    # max_num = 0
    # for s in slice_final:
    #     # print(s[0])
    #     if s[0] > max_num:
    #         max_num = s[0]
    # print(f"{time_interval}  {machine_id}机床的运行状态切片最大值为：", max_num)
    publish_message(LOW_FREQ, j)

    return temp_dict


def UpdMachineEnergy(time_interval: str) -> None:
    """2-16 更新所有机床能耗  车间大屏 - 机床能耗排行【假数据】"""
    # dice = random.randint(0, 1)
    # rank = [[["H6000", "H8000", "NHM6300"], [89, 86, 90]], [["H6000", "H8000", "NHM6300"], [88, 91, 87]]]
    temp_dict = {}
    temp_dict["task"] = str("update_machine_energy")
    try:
        temp_dict["time_interval"] = str(time_interval)
        # temp_dict["energy_rank"] = rank[0]
        temp_dict["energy_rank"] = [["H6000", "H8000", "NHM6300"], [89, 86, 90]]
        j = json.dumps(temp_dict, ensure_ascii=False)
    except KeyError:
        temp_dict["energy_rank"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
        j = json.dumps(temp_dict, ensure_ascii=False)
    publish_message(LOW_FREQ, j)


def UpdEnergyPartion(time_interval: str) -> None:
    """2-18 更新总体能耗构成  【更新大屏前端后此接口作废】"""
    # dice = random.randint(0, 1)
    # energy_partion = [{"running": 68, "shutdown": 23, "breakdown": 7, "pausing": 2},
    #                   {"running": 89, "shutdown": 1, "breakdown": 1, "pausing": 9}]
    temp_dict = {}
    temp_dict["task"] = str("update_energy_partion")
    try:
        temp_dict["time_interval"] = str(time_interval)
        # TODO: 最新版 修改完待测试
        temp_dict["energy_partion"] = {'running': 32, 'shutdown': 23, 'breakdown': 2, 'pausing': 43}
        j = json.dumps(temp_dict, ensure_ascii=False)
    except KeyError:
        temp_dict["energy_partion"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
        j = json.dumps(temp_dict, ensure_ascii=False)
    publish_message(LOW_FREQ, j)


def UpdEnergyPartionMachine(time_interval: str, machine_id: int) -> None:
    """
    2-19 更新单台机床能耗  【更新大屏前端后此接口作废】
    """
    # dice = random.randint(0, 1)
    # energy_partion = [{"running": 68, "shutdown": 23, "breakdown": 7, "pausing": 2},
    #                   {"running": 89, "shutdown": 1, "breakdown": 1, "pausing": 9}]
    temp_dict = {}
    temp_dict["task"] = str("update_energy_partion_machine")
    try:
        temp_dict["time_interval"] = str(time_interval)
        temp_dict["machine_id"] = str(machine_id)
        # TODO: 最新版 修改完待测试
        temp_dict["energy_partion"] = {'running': 32, 'shutdown': 23, 'breakdown': 2, 'pausing': 43}
        j = json.dumps(temp_dict, ensure_ascii=False)
    except KeyError:
        temp_dict["energy_partion"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
        j = json.dumps(temp_dict, ensure_ascii=False)
    publish_message(LOW_FREQ, j)


def UpdCarbonMachine(time_interval: str, machine_id: int) -> None:
    """ 2-20 车间总览 - 碳排放【假数据】  【更新大屏前端后此接口作废】"""
    temp_dict = {}
    temp_dict["task"] = str("update_workshop_carbon_emission_machine")
    # return None
    try:
        temp_dict["time_interval"] = str(time_interval)
        temp_dict["machine_id"] = str(machine_id)
        # temp_dict["carbon_emission"] = carbon_emissions[dice][time_interval]
        # temp_dict["carbon_per_money"] = carbon_per_money[dice][time_interval]
        temp_dict["carbon_emission"] = ["200.0", "100.0", "200.0", "343.0", "250.0", "112.0"]
        temp_dict["carbon_per_money"] = ["200.0", "100.0", "200.0", "343.0", "250.0", "112.0"]
        j = json.dumps(temp_dict, ensure_ascii=False)
    except KeyError:
        temp_dict["carbon_emission"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
        j = json.dumps(temp_dict, ensure_ascii=False)
    # print("carbon_emissions: ", carbon_emissions)
    # print("carbon_per_money: ", carbon_per_money)
    # print("更新单台机床碳排放: ", machine_id)
    # print("数据：", j)
    publish_message(LOW_FREQ, j)


def UpdOEEDaily(time_interval: str, machine_id: int):
    """2-24 机床状态 - 设备分时有效运行率"""
    # TODO： target : dict{str(日期),dict{str(状态),float(数值)}}    current : dict{str(状态),list(数值列表)}
    temp_dict = {}
    temp_dict["task"] = str("update_machine_OEE_daily")
    temp_dict["time_interval"] = str(time_interval)
    temp_dict["machine_id"] = str(machine_id)
    start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
    OEE_list = _update_machine_slice_daily(machine_id, start_time, end_time)
    if not OEE_list:
        return None
    OEE_final = {}
    for key, value in OEE_list.items():
        OEE_final[MACHINE_STATUS_INV[str(key)]] = value

    # # # print(OEE_final)
    try:
        temp_dict["OEE_daily"] = OEE_final
        j = json.dumps(temp_dict, ensure_ascii=False)
    except KeyError:
        temp_dict["OEE_daily"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
        j = json.dumps(temp_dict, ensure_ascii=False)
    # print(f"更新设备分时OEE: ", j)
    publish_message(LOW_FREQ, j)

    return temp_dict


def UpdProgTimes(time_interval: str) -> None:
    """2-25 更新所有机床累计运行任务数  车间大屏 - 运行任务"""
    machine_num = __get_machine_num()
    machine_id_list = __get_machine_id_list()
    start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
    total_prog_number = 0
    for i in range(machine_num):
        machine_prog_rank = _update_machine_prog_rank(
            str(machine_id_list[i]),
            start_time,
            end_time,
        )

        if machine_prog_rank is None:
            return None

        prog_number, run_number = machine_prog_rank

        total_prog_number += prog_number

    temp_dict = {}
    temp_dict["task"] = str("update_all_prog_times")
    temp_dict["time_interval"] = str(time_interval)
    try:
        temp_dict["times"] = total_prog_number
        j = json.dumps(temp_dict, ensure_ascii=False)
    except KeyError:
        temp_dict["times"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
        j = json.dumps(temp_dict, ensure_ascii=False)
    publish_message(LOW_FREQ, j)

    return temp_dict


# 更新机床利用率
# FIXME 改成OEE的数据
# def UpdMachineUseRate(time_interval: str, machine_id: int) -> None:
#     dice = random.randint(0, 1)
#     userate = [
#         {
#             "day-4hour": {"running": 0.6, "shutdown": 0.2, "breakdown": 0.1, "pausing": 0.1},
#             "week-1day": {"running": 0.6, "shutdown": 0.2, "breakdown": 0.1, "pausing": 0.1},
#             "month-5day": {"running": 0.6, "shutdown": 0.2, "breakdown": 0.1, "pausing": 0.1},
#         },
#         {
#             "day-4hour": {"running": 0.2, "shutdown": 0.6, "breakdown": 0.1, "pausing": 0.1},
#             "week-1day": {"running": 0.1, "shutdown": 0.7, "breakdown": 0.1, "pausing": 0.1},
#             "month-5day": {"running": 0.6, "shutdown": 0.1, "breakdown": 0.1, "pausing": 0.2},
#         },
#     ]
#     temp_dict = {}
#     temp_dict["task"] = str("update_machine_use_rate")
#     try:
#         temp_dict["time_interval"] = str(time_interval)
#         temp_dict["machine_id"] = str(machine_id)
#         temp_dict["machine_use_rate"] = userate[dice][time_interval]
#         j = json.dumps(temp_dict, ensure_ascii=False)
#     except KeyError:
#         temp_dict["times"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
#         j = json.dumps(temp_dict, ensure_ascii=False)


#     print("机床利用率: ",j)
#     publish_message(LOW_FREQ, j)


def UpdMachineUseRate(time_interval: str, machine_id: int):
    """2-26 车间大屏 - 左边的设备有效运行率"""

    start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
    try:
        oee_list = _update_machine_OEE(machine_id=machine_id, start_time=start_time, end_time=end_time)
    except Exception as e:
        print(e)
    if not oee_list:
        return None
    # print(oee_list)

    idle = float(str(oee_list["IDLE"][1]).replace("%", "")) / 100
    not_connect = float(str(oee_list["NOT_CONNECT"][1]).replace("%", "")) / 100
    pausing = float(str(oee_list["PAUSING"][1]).replace("%", "")) / 100
    running = float(str(oee_list["RUNNING"][1]).replace("%", "")) / 100

    temp_dict = {
        "task": str("update_machine_use_rate"),
        "machine_id": str(machine_id),
        "machine_use_rate": {
            # 保留小数点后1位 更改为 保留小数点后2位
            "running": round(running, 2),
            # "shutdown": round(idle + not_connect, 2),
            # "breakdown": 0,
            # TODO: 最新版 修改逻辑，待测试
            "shutdown": round(idle, 2),
            "breakdown": round(not_connect, 2),
            "pausing": round(pausing, 2),
        },
    }
    try:
        temp_dict["time_interval"] = str(time_interval)
    except KeyError:
        temp_dict["times"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
    # print(f"机床利用率：  ",temp_dict)
    publish_message(LOW_FREQ, json.dumps(temp_dict, ensure_ascii=False))


def UpdMachineOEE(UUID: str, machine_id: int, time_interval: str):
    """
    3-98 点击机床（没点进入） - 显示的OEE
    进入机床 - 显示的OEE
    机床动态 - 设备总体有效运行率
    """

    temp_dict = {}
    temp_dict["task"] = str("ask_for_machine_OEE")
    temp_dict["UUID"] = str(UUID)
    temp_dict["time_interval"] = str(time_interval)
    temp_dict["machine_id"] = str(machine_id)
    start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
    status_list = ["running", "breakdown", "pausing", "shutdown"]
    dict_list = [0, 0, 0, 0]
    OEE_list = _update_machine_OEE(machine_id, start_time, end_time)

    if OEE_list is None:
        return None

    for j in range(len(status_list)):
        dict_list[j] = OEE_list[MACHINE_STATUS[status_list[j]]][0]
        # TODO: 添加 保留小数点后2位 操作
    # print(f"单个机床的OEE: ", dict_list)
    per_list = [str(round(x / sum(dict_list), 2)) for x in dict_list]

    # 7-1 防止per_list列表总和小于1
    per_list_int = [float(x) for x in per_list]
    if sum(per_list_int) < 1:
        difference = 1 - sum(per_list_int)
        # 如果总和少于1，则将缺少的部分填充到running状态中，使其的总和保证 = 1
        per_list[0] = str(float(per_list[0]) + difference)

    # 防止per_list列表总和大于1
    if sum(per_list_int) > 1:
        difference = sum(per_list_int) - 1
        # 如果总和大于1，则将多余的部分从pausing状态中删除，使其的总和保证 = 1
        per_list[2] = str(float(per_list[2]) - difference)

    # 7-2 防止出现有些状态的百分比大于100，有些百分比小于0
    for per in per_list:
        if float(per) > 1 or float(per) < 0:
            per_list = ["0.4", "0.0", "0.6", "0.0"]
            break

    try:
        temp_dict["pie_list"] = [status_list, per_list]
        j = json.dumps(temp_dict, ensure_ascii=False)
    except KeyError:
        temp_dict["pie_list"] = "时间间隔不正确, 请保证时间间隔是day-4hour,week-1day,month-5day之一"
        j = json.dumps(temp_dict, ensure_ascii=False)
    # print(f"更新单个机床OEE: ", j)
    publish_message(LOW_FREQ, j)

    return temp_dict


# FIXME 获取实时数据
def UpdMachine(UUID: str, machine_id: int):
    """
    3-99 更新机床静态数据
    走到机床前显示的静态数据
    点击机床后显示的静态数据
    机床动态 - 机床参数
    """
    # MACHINES_CONFIG = os.path.join(str(path_depend), MACHINE_CONFIG)
    # api_document = parse(get_conf_file("API_config.xml"))
    # api = api_document.getElementsByTagName("API")

    # machines_config = get_conf_file(MACHINE_CONFIG)
    # domTree = parse(machines_config)
    # rootNode = domTree.documentElement
    # machines = rootNode.getElementsByTagName("Machine")
    temp_dict = {}
    global OLD_DATA
    global TOOL
    # connector = RedisConnector()

    for i in machine:
        if i.getElementsByTagName("id")[0].childNodes[0].data == str(machine_id):
            machine_select = i
            static = machine_select.getElementsByTagName("Static")[0]
            basic = machine_select.getElementsByTagName("Basic")[0]
            # high_freq_element = machine_select.getElementsByTagName("HighFreq")[0]
            machine_name = static.getElementsByTagName("machine_name")[0].childNodes[0].data
            temp_dict = {}
            temp_dict["task"] = str("ask_for_machine")
            temp_dict["UUID"] = str(UUID)
            temp_dict["machine_id"] = str(machine_id)
            temp_dict["machine_type"] = str(
                static.getElementsByTagName("machine_type")[0].childNodes[0].data,
            )
            temp_dict["machine_name"] = str(
                static.getElementsByTagName("machine_name")[0].childNodes[0].data,
            )

            temp_dict["machine_merchant"] = str(
                static.getElementsByTagName("machine_merchant")[0].childNodes[0].data,
            )
            temp_dict["text_id"] = str(static.getElementsByTagName("text_id")[0].childNodes[0].data)
            temp_dict["power(kW)"] = str(
                static.getElementsByTagName("machine_power")[0].childNodes[0].data,
            )
            temp_dict["nc_system"] = str(
                static.getElementsByTagName("machine_nc")[0].childNodes[0].data,
            )

            if machine_name != "NHM6300":
                # data = connector.get_data_from_single_machine(machine_index=str(machine_id))
                host = basic.getElementsByTagName("opc_host")[0].childNodes[0].data
                port = basic.getElementsByTagName("opc_port")[0].childNodes[0].data
                username = basic.getElementsByTagName("opc_username")[0].childNodes[0].data
                password = basic.getElementsByTagName("opc_password")[0].childNodes[0].data

                # 实时数据
                no_data = ""
                name_node = "ns=2;s=/Channel/ProgramInfo/selectedWorkPProg"
                name_data = get_tool_compensation_to_opc(name_node, host, port, username, password)
                # print("+++/////////////////+++name_data", name_data)
                # 防止name_data为空获取不到运行程序名
                if name_data is None:
                    name_data = OLD_DATA
                else:
                    OLD_DATA = name_data
                running_program = reduce_program_name(
                    str(name_data))
                # print("--------------///////////////running_program", running_program)
                workpiece_info = workpiece.get(running_program)
                #                 print("workpiece_info", workpiece_info)
                if workpiece_info is None:
                    temp_dict["product"] = no_data
                    temp_dict["product_id"] = no_data
                else:
                    product = workpiece[running_program].get("product")
                    pid = workpiece[running_program].get("product_id")
                    # print("///////////product", product)
                    # print("//////////pid", pid)
                    if product is None:
                        temp_dict["product"] = no_data
                    else:
                        temp_dict["product"] = product

                    if pid is None:
                        temp_dict["product_id"] = no_data
                    else:
                        temp_dict["product_id"] = pid

                # 运行程序
                temp_dict["running_program"] = running_program
                # 刀具号
                tool_node = "ns=2;s=/Channel/State/actTNumber"
                try:
                    tool_data = get_tool_compensation_to_opc(tool_node, host, port, username, password)
                    # print("tool_data", tool_data)
                    if tool_data is None:
                        temp_dict["tool"] = "暂无"
                        # temp_dict["tool"] = str(tool_list[random.randint(0, 9)])
                    else:
                        temp_dict["tool"] = tool_data
                        # rows = select_from_single_machine(
                        #     machine_index=machine_id,
                        #     table_name="tool_mapping",
                        #     column_name=["tool_id"],
                        #     where=f"tool_no = {tool_data}",
                        # )
                        # if rows:
                        #     temp_dict["tool"] = rows[0][0]
                except:
                    temp_dict["tool"] = "暂无"
                    # temp_dict["tool"] = str(tool_list[random.randint(0, 9)])
                TOOL = {"tool": temp_dict["tool"]}
                # temp_dict["fault_code"] = "404"
                # temp_dict["fault_information"] = "机床未连接"
                # j = json.dumps(temp_dict, ensure_ascii=False)
                # publish_message(LOW_FREQ, j)
                # print("-----________________------发送redis成功了!!!!")
            else:  # TODO: 按照machine_name != "NHM6300"时的逻辑获取product和product_id数据
                try:
                    connector = RedisConnector()
                    data = connector.get_data_from_single_machine(machine_id)
                except:
                    print("》》》》》》》》》》》》》》》》》》》》》》》》》》NHM6300的redis连接出错")
                    return
                # temp_dict["product"] = str(product_list[random.randint(0, 2)])
                # temp_dict["product_id"] = str(product_id_list[random.randint(0, 2)])
                # temp_dict["running_program"] = str("0002.SPF")
                # temp_dict["tool"] = str("80")
                # temp_dict["fault_code"] = "404"
                # temp_dict["fault_information"] = "机床未连接"
                temp_dict["running_program"] = data['name']
                temp_dict["tool"] = data['tool']
                no_data = ""
                workpiece_info = workpiece.get(data['name'])
                # print("workpiece_info", workpiece_info)
                if workpiece_info is None:
                    temp_dict["product"] = no_data
                    temp_dict["product_id"] = no_data
                else:
                    product = workpiece[data['name']].get("product")
                    pid = workpiece[data['name']].get("product_id")
                    # print("product", product)
                    # print("pid", pid)
                    if product is None:
                        temp_dict["product"] = no_data
                    else:
                        temp_dict["product"] = product

                    if pid is None:
                        temp_dict["product_id"] = no_data
                    else:
                        temp_dict["product_id"] = pid
                # print("------------------------------NHM6300的数据： ", temp_dict)
            j = json.dumps(temp_dict, ensure_ascii=False)
            publish_message(LOW_FREQ, j)
    return temp_dict
