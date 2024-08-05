from __future__ import unicode_literals

import json
import os
from typing import Dict
from xml.dom.minidom import parse

import pandas as pd
import requests

from lib.Config.config import MACHINE_CONFIG, STATIC_INFO, get_conf_file
from lib.Interface.mysqlConnector import select_from_single_machine
from lib.Interface.redis_api import publish_message
from .lowFreqFresh import TOOL, __gen_datetime_interval

# path = os.path.abspath(os.path.dirname(__file__))
# FILE = os.path.join(path, "backend.xlsx")

# FILE = get_conf_file("backend.xlsx")
### 建立mock表
# machine_list = [  # False
#     "NB251",
#     "NB252",
#     "STC600",
#     "STC601",
#     "LX051",
#     "NB251",
#     "NB252",
#     "STC600",
#     "STC601",
#     "LX051",
# ]
# machine_list = pd.read_csv(os.path.join('conf', 'mock_machine_list.csv'))['machine_list'].tolist()

# machine_type_list = ["五轴加工中心", "车铣复合中心", "车床", "磨床", "钻削中心", "龙门式车床"]  # False
# machine_type_list = pd.read_csv(os.path.join('conf', 'mock_machine_type.csv'))['machine_type'].tolist()
# print('machine_type_list: ', machine_type_list)

# machine_status = ["running", "shutdown", "breakdown", "pausing"]  # False
# machine_status = pd.read_csv(os.path.join('conf', 'mock_machine_status.csv'))['machine_status'].tolist()
# print('machine_status: ', machine_status)

# tool_list = ["12", "34", "67", "15", "10", "76", "59", "23", "61", "316"]  # False
# tool_list = pd.read_csv(os.path.join('conf', 'mock_tool_list.csv'))['tool_list'].tolist()
# print('tool_list: ', tool_list)

# program_list = [  # False
#     "/_N_WKS_DIR/_N_TEMP_WPD/_N_Z130_2643K2_02H_ZAZK_R_X_MPF",
#     "/_N_WKS_DIR/_N_DWED_WPD/_N_Z130_2643K2_02H_ZAZK_R_X_MPF",
#     "/_N_WKS_DIR/_N_TEMP_WPD/_N_Z130_2DWEDW2_02H_ZAZK_R_X_MPF",
#     "/_N_WKS_DIR/_N_TEMP_WPD/_N_Z130_2643K2_02H_ZAZK_R_X_MPF",
#     "/_N_WKS_DIR/_N_TEDWEP_WPD/_N_Z130_2643K2_02H_ZAZK_R_X_MPF",
#     "/_N_WKS_DIR/_N_TEMP_WPD/_N_Z130_264WEDW2_02H_ZAZK_R_X_MPF",
#     "/_N_WKS_DIR/_N_TEMP_WPD/_N_Z130_2EDWK2_02H_ZAZK_R_X_MPF",
#     "/_N_WKS_DIR/_N_TEMP_WPD/_N_Z130_2643K2_02H_ZADWEDEW_R_X_MPF",
#     "/_N_WWDEDWS_DIR/_N_TEMP_WPD/_N_Z130_2643K2_02H_ZAZK_R_X_MPF",
#     "/_N_WKS_DIR/_N_TEMP_WPD/_N_Z1WDE_2643K2_02H_ZAZK_R_X_MPF",
# ]
# program_list = pd.read_csv(os.path.join('conf', 'mock_program_list.csv'))['program_list'].tolist()
# print('program_list: ', program_list)

# time_interval = ["day-2hour", "week-1day", "month-3day"]  # False
# time_interval = pd.read_csv(os.path.join('conf', 'mock_time_interval.csv'))['time_interval'].tolist()
# print('time_interval: ', time_interval)

# machine_basic_information = {
#     "主轴接口": 3,
#     "支持轴数": 4,
#     "数控系统": "Simens840D",
#     "主轴最大转速": 8000,
#     "额定功率": 40,
#     "入厂时间": "2022-05-01",
# }
# machine_constraints = {
#     "最大工件直径": 1290,
#     "最大工件高度": 1200,
#     "最大工件重量": 1400,
# }
# machine_tool_magazine = {
#     "刀库类型": "ATC 刀库",
#     "刀库容量": 100,
#     "刀位允许直径": 112,
#     "刀位最大承重": 35,
#     "刀位允许长度": 600,
# }
# machine_cool = {
#     "支持冷却类型": "外冷/内冷",
#     "支持冷却液类型": "5%乳化液",
#     "内冷却压力": 50,
# }
# machine_precision = {
#     "定位精度": 0.006,
#     "重复定位精度": 0.0024,
#     "反向差值": 0.002,
#     "回转定位精度": 7,
#     "回转重复定位精度": 5,
# }
# machine_power = {
#     "主轴功率曲线": [[1, 2, 3, 4, 5, 6, 7, 8, 9], [4000, 4000, 4000, 5000, 6000, 10000, 10000, 15000, 20000]],
# }
# machine_torque = {
#     "主轴扭矩曲线": [[1, 2, 3, 4, 5, 6, 7, 8, 9], [10000, 10000, 10000, 15000, 10000, 10000, 10000, 15000, 20000]],
# }
# machine_health = {
#     "机床体检评分": [[1, 2, 3, 4, 5, 6, 7, 8, 9], [100, 100, 100, 100, 100, 100, 100, 100, 100]],
# }
machine_specifications = pd.read_csv(os.path.join('conf', 'mock_machine_specifications.csv'))
# machine_data = machine_specifications[machine_specifications['id'] == 801022]
# # print("1: ", machine_data)
# # 确保只有一条记录匹配machine_id（假设ID唯一）
# machine_data = machine_data.iloc[0]
# # print("2: ", machine_data)
# machine_basic_information = machine_specifications.iloc[0, 1:7].to_dict()
# # print(machine_basic_information)
#
# machine_constraints = machine_specifications.iloc[0, 7:10].to_dict()
# # print(machine_constraints)
#
# machine_tool_magazine = machine_specifications.iloc[0, 10:15].to_dict()
# # print(machine_tool_magazine)
#
# machine_cool = machine_specifications.iloc[0, 15:18].to_dict()
# # print(machine_cool)
#
# machine_precision = machine_specifications.iloc[0, 18:23].to_dict()
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
# # point = ["2022-06-01", "2022-06-02", "2022-06-03", "2022-06-04", "2022-06-05", "2022-06-06", "2022-06-07", "2022-06-08",
# #          "2022-06-09"]
# point = machine_specifications["机床体检日期"].tolist()[1:]
# machine_health["机床体检评分"] = [point, machine_health_curve]
# # print(machine_health)

### 默认图形库
# pie_list = [  # False
#     ["0.1", "0.2", "0.4", "0.1"],
#     ["0.2", "0.4", "0.3", "0.1"],
#     ["0.2", "0.5", "0.2", "0.1"],
#     ["0.3", "0.4", "0.2,", "0.1"],
#     ["0.2", "0.7", "0.0", "0.1"],
# ]
# pie_list = [
#     [str(value) for value in row]
#     for row in pd.read_csv(os.path.join('conf', 'mock_pie_data.csv')).values.tolist()
# ]
# print('pie_list: ', pie_list)

# finish_list = [
#     ["43", "32", "21"],
#     ["646", "544", "124"],
#     ["44", "1", "0"],
#     ["65", "55", "3"],
#     ["21", "3", "1"],
# ]
# finish_rate_list = [  # False
#     ["0.7", "0.5", "0.5"],
#     ["0.6", "0.4", "0.2"],
#     ["0.9", "0.5", "0.3"],
#     ["0.9", "0.4", "0.3"],
#     ["0.8", "0.7", "0.1"],
# ]
# finish_rate_list = [
#     [str(value) for value in row]
#     for row in pd.read_csv(os.path.join('conf', 'mock_finish_rate_list.csv')).values.tolist()
# ]
# print(finish_rate_list)
# list_seven = [  # False
#     ["200.0", "100,0", "200,0", "343.0", "250.0", "112.0", "234.9"],
#     ["230.0", "400,0", "240,0", "321.0", "243.0", "156.0", "454.9"],
#     ["243.0", "123,0", "780,0", "398.0", "240.0", "113.0", "124.9"],
#     ["420.0", "130,0", "300,0", "300.0", "980.0", "123.0", "254.9"],
#     ["670.0", "300,0", "500,0", "300.0", "230.0", "145.0", "124.9"],
# ]
# list_seven_list = [
#     [str(value) for value in row]
#     for row in pd.read_csv(os.path.join('conf', 'mock_list_seven.csv')).values.tolist()
# ]
# print(list_seven_list)

# product_list = ["机匣", "细长轴", "轴承外圈"]

# product_list: list[str] = ["连杆", "连杆", "连杆"]

# product_id = ["WX15736372", "STE78382782", "OIU-9032390"]

# product_id_list = ["091123", "324242", "324423"]  # False
# product_id_list = [  # TODO: 想要获取"091123"，获取到"91123"，会自动将前导0去除
#     [str(value) for value in row]
#     for row in pd.read_csv(os.path.join('conf', 'mock_product.csv')).values.tolist()
# ]
# print(product_id_list)

# data = pd.read_excel(FILE, header=None, engine="openpyxl")
# data = pd.DataFrame(data.values.T, columns=data.values.T[0, :])
# data = data.drop(index=0)

# tool_info = {
#     "105": {
#         "type": "倒角",
#         "blade_number": 3,
#     },
#     "703": {
#         "type": "合金钻头",
#         "blade_number": 2,
#     },
#     "402": {
#         "type": "面铣刀",
#         "blade_number": 8,
#     },
#     "D034": {
#         "type": "U钻",
#         "blade_number": 2,
#     },
#     "D38": {
#         "type": "U钻",
#         "blade_number": 2,
#     },
#     "403": {
#         "type": "面铣刀",
#         "blade_number": 8,
#     },
# }
# machine_data = machine_specifications[machine_specifications['id'] == 801015]
# machine_data = machine_data.iloc[0]
#
# machine_basic_information = machine_data[['主轴接口', '支持轴数', '数控系统', '主轴最大转速', '额定功率', '入厂时间']].to_dict()
# print(machine_basic_information)
# machine_constraints = machine_data[['最大工件直径', '最大工件高度', '最大工件重量']].to_dict()
# print(machine_constraints)
# machine_tool_magazine = machine_data[['刀库类型', '刀库容量', '刀位允许直径', '刀位最大承重', '刀位允许长度']].to_dict()
# print(machine_tool_magazine)
# machine_cool = machine_data[['支持冷却类型', '支持冷却液类型', '内冷却压力']].to_dict()
# print(machine_cool)
# machine_precision = machine_data[['定位精度', '重复定位精度', '反向差值', '回转定位精度', '回转重复定位精度']].to_dict()
# print(machine_precision)

df_tools = pd.read_csv(os.path.join('conf', 'mock_tool_info.csv'))
# 构建 tool_info 字典
tool_info = {
    tool_id: {"type": tool_type, "blade_number": blade_number}
    for tool_id, tool_type, blade_number in zip(
        df_tools["id"], df_tools["type"], df_tools["blade_number"]
    )
}

# 临时工件信息
# workpiece: Dict[str, Dict[str, str]] = {
#     "O2760": {
#         # 零件名称
#         "product": "L27连杆-80",
#         # 零件型号
#         "product_id": "Z11.03000-0822-1.I.a",
#         # 零件数
#         "finished_number": "0",
#         # 零件信息
#         "product_remain": "加工杆身四周面",
#     },
#     "O2761": {
#         # 零件名称
#         "product": "L27连杆-90-2-l",
#         # 零件型号
#         "product_id": "Z11.03000-0822-1.K.a",
#         # 零件数
#         "finished_number": "0",
#         # 零件信息
#         "product_remain": "卧加杆盖, 中段外侧三面",
#     },
#     "O2762": {
#         # 零件名称
#         "product": "L27连杆-90-2-r",
#         # 零件型号
#         "product_id": "Z11.03000-0822-1.K.a",
#         # 零件数
#         "finished_number": "0",
#         # 零件信息
#         "product_remain": "卧加杆盖, 中段外侧三面",
#     },
#     "O2763": {
#         # 零件名称
#         "product": "L27连杆-90-3-l",
#         # 零件型号
#         "product_id": "Z11.03000-0822-1.K.a",
#         # 零件数
#         "finished_number": "0",
#         # 零件信息
#         "product_remain": "卧加杆盖, 中段哈佛面",
#     },
#     "O2764": {
#         # 零件名称
#         "product": "L27连杆-90-3-r",
#         # 零件型号
#         "product_id": "Z11.03000-0822-1.K.a",
#         # 零件数
#         "finished_number": "0",
#         # 零件信息
#         "product_remain": "卧加杆盖, 中段哈佛面",
#     },
# }
# df_workpiece = pd.read_csv(os.path.join('conf', 'mock_workpiece.csv'))
# # 构建 workpiece 字典
# workpiece = {
#     str(workpiece_id): {"product": product, "product_id": product_id, "finished_number": finished_number,
#                         "product_remain": product_remain}
#     for workpiece_id, product, product_id, finished_number, product_remain in zip(
#         df_workpiece["id"], df_workpiece["product"], df_workpiece["product_id"], df_workpiece["finished_number"],
#         df_workpiece["product_remain"]
#     )
# }
# print(workpiece)
machines_config = get_conf_file(MACHINE_CONFIG)
domTree = parse(machines_config)
rootNode = domTree.documentElement
machine = rootNode.getElementsByTagName("Machine")

df_workpiece = pd.read_csv(os.path.join('conf', 'mock_workpiece.csv'))
# 构建 workpiece 字典
workpiece = {
    str(product): {"product_id": product_id, "finished_number": finished_number,
                   "product_remain": product_remain}
    for product, product_id, finished_number, product_remain in zip(
        df_workpiece["product"], df_workpiece["product_id"], df_workpiece["finished_number"],
        df_workpiece["product_remain"]
    )
}


# def UpdMachineTool(tool: str) -> None:
#     """
#     4-1 获取刀具信息
#     TODO 接口调整, 添加新的入参：机床编号(machine_id)
#     1. 根据机床编号区分excel表
#     2. 根据刀具名称从excel表中查询刀具信息
#     """
#
#     temp_dict = {}
#     temp_dict["task"] = str("update_machine_tool")
#     temp_dict["tool"] = str(tool)
#     detail = ""

# for i in machine:
#     machine_id = i.getElementsByTagName("id")[0].childNodes[0].data
#     rows = select_from_single_machine(
#         machine_index=machine_id,
#         table_name="tool_mapping",
#         column_name=["edge_length", "edge_radius"],
#         where=f"tool_id = '{tool}'",
#     )
#     if rows:
#         edge_length, edge_radius = rows[0]
#         detail = f"{detail};刀补长度: {edge_length};刀补半径: {edge_radius}"
#         break

# ti = tool_info.get(tool)
# if ti:
#     tool_type = ti.get("type")
#     blade_number = ti.get("blade_number")
#     if tool_type:
#         detail = str(tool_type)
#     if blade_number:
#         detail = f"{detail};刃数: {blade_number}"
#
# temp_dict["tool_merchant"] = str("sandvik")
# temp_dict["tool_remain"] = detail
#
# j = json.dumps(temp_dict, ensure_ascii=False)
# publish_message(STATIC_INFO, j)


#
# def UpdMachineProduct(product: str) -> None:
#     """
#     4-2 获取零件信息
#     TODO 接口需要调整, 添加新的入参: 机床编号(machine_id)、当前正在运行的程序名称(program_running)
#     1. 根据机床编号和程序名称查询对应program表, 计算加工完成的零件数量
#     2. 根据零件名称和程序名称在excel表中(前期, 后期从MES系统接口中获取)查找对应的零件信息
#     """
#
#     # print("product: ", product)
#
#     temp_dict = {}
#     program_name = ""
#     product_id = ""
#     finished_number = ""
#     product_remain = ""
#
#     for pro_name, workpiece_info in workpiece.items():
#         p = workpiece_info.get("product")
#         if p == product:
#             program_name = pro_name
#             product_id = workpiece_info.get("product_id")
#             product_remain = workpiece_info.get("product_remain")
#             break
#
#     if program_name:
#         # 如果程序名称存在
#
#         finish_list = []
#         for i in machine:
#             machine_id = i.getElementsByTagName("id")[0].childNodes[0].data
#             rows = select_from_single_machine(
#                 machine_index=machine_id,
#                 table_name="program",
#                 column_name=["COUNT(id)"],
#                 where=f"`name` LIKE '%{program_name}%' AND end_time IS NOT NULL AND run_time > 30",
#             )
#             if rows:
#                 finish_list.append(str(rows[0][0]))
#
#         if finish_list:
#             temp_dict["finished_number"] = str(max(finish_list))
#         else:
#             temp_dict["finished_number"] = finished_number
#
#     temp_dict["task"] = str("update_machine_product")
#     temp_dict["product"] = product.split("-")[0]
#     temp_dict["product_id"] = product_id
#     temp_dict["product_remain"] = product_remain
#
#     publish_message(STATIC_INFO, json.dumps(temp_dict, ensure_ascii=False))


# 4-3

# 函数为点击触发（点击进入机床）
def UpdMachineTool(tool: str) -> None:
    temp_dict = {}
    temp_dict["task"] = str("update_machine_tool")
    if tool != "":
        temp_dict["tool"] = str(tool)
    else:
        if TOOL != {}:
            temp_dict["tool"] = TOOL["tool"]
        else:
            temp_dict["tool"] = "暂无"

    detail = "暂无"
    ti = tool_info.get(tool)
    if ti:
        tool_type = ti.get("type")
        blade_number = ti.get("blade_number")
        if tool_type:
            detail = str(tool_type)
        if blade_number:
            detail = f"{detail};刃数: {blade_number}"
    temp_dict["tool_remain"] = str(detail)
    # TODO: 最新版 修改完待测试
    # temp_dict["tool_merchant"] = str("sandvik")
    j = json.dumps(temp_dict, ensure_ascii=False)
    print("updMachineTool   ", j)
    publish_message(STATIC_INFO, j)


# 函数为点击触发（点击进入机床）
def UpdMachineProduct(product: str, machine_id: str) -> None:
    temp_dict = {}
    temp_dict["task"] = str("update_machine_product")

    # TODO： 获取机床的零件完成数据(固定获取一个月内mes中已完成数的数据)
    workpiece_finished_number = "0"
    for i in machine:
        if i.getElementsByTagName("id")[0].childNodes[0].data == str(machine_id):
            machine_select = i
            static = machine_select.getElementsByTagName("Static")[0]
            machine_name = static.getElementsByTagName("machine_name")[0].childNodes[0].data

            time_interval = "month-5day"
            start_time, end_time = __gen_datetime_interval(time_interval=time_interval)
            # WORKPIECE_FINISHED_COUNT = __select_program_count_from_machine(machine_id, start_time, end_time)
            # print(f"{machine_id}机床完成零件数", WORKPIECE_FINISHED_COUNT)

            url = "http://192.168.25.140:8081/mes-web/zcplanmanagerController/zyfhGxrw11.do"
            if machine_name == "NHM6300":
                # 查询字符串参数
                params = {
                    "manuCode": "H6300",
                    "startDate": start_time.strftime("%Y-%m-%d"),
                    "endDate": end_time.strftime("%Y-%m-%d")
                }
            else:
                params = {
                    "manuCode": machine_name,
                    "startDate": start_time.strftime("%Y-%m-%d"),
                    "endDate": end_time.strftime("%Y-%m-%d")
                }
            # print(f"请求的params： ", params)
            response = requests.get(url, params=params)
            # print("响应状态码：", response.status_code)
            # 检查响应状态码
            if response.status_code == 200:
                res = response.json()
                # print(f"请求的机床{machine_id}, 获取的数据为：{res}")
                if res:
                    tmp = res[0]['SHULIANG']
                    workpiece_finished_number = tmp.split("/")[0]
    if product == "":
        temp_dict["product"] = "暂无"
        temp_dict["product_id"] = "暂无"
        temp_dict["finished_number"] = workpiece_finished_number
        temp_dict["product_remain"] = "暂无"
    else:
        temp_dict["product"] = product  # 加工零件名称
        # temp_dict["product_id"] = str("484646")  # 加工零件型号
        # temp_dict["finished_number"] = str("213")  # 加工零件数
        # temp_dict["product_remain"] = str("保留信息")
        # 根据零件名称在program表中找到end_time is not null AND 一段时间间隔内的记录 AND name == product的记录数量（成功加工的product零件数量）
        temp_dict["product_id"] = str(workpiece[product]["product_id"])
        temp_dict["product_remain"] = str(workpiece[product]["product_remain"])
        # temp_dict["finished_number"] = str(workpiece[product]["finished_number"])
        temp_dict["finished_number"] = workpiece_finished_number
    j = json.dumps(temp_dict, ensure_ascii=False)
    print("updMachineProduct   ", j)
    publish_message(STATIC_INFO, j)


def UpdMachineInfo(machine_id: str) -> None:
    """ 更新生产资料 """
    # TODO: 6-17 修改函数
    temp_dict = {}

    machine_data = machine_specifications[machine_specifications['id'] == machine_id]

    if machine_data.empty:
        print(f"No data found for machine_id: {machine_id}")
        return

    # 确保只有一条记录匹配machine_id（假设ID唯一）
    machine_data = machine_data.iloc[0]

    machine_basic_information = machine_data[
        ['主轴接口(个)', '支持轴数(个)', '数控系统', '主轴最大转速(rpm)', '额定功率(w)', '入厂时间']].to_dict()
    # print(machine_basic_information)

    machine_constraints = machine_data[['最大工件直径(mm)', '最大工件高度(mm)', '最大工件重量(g)']].to_dict()
    # print(machine_constraints)

    machine_tool_magazine = machine_data[
        ['刀库类型', '刀库容量(个)', '刀位允许直径(mm)', '刀位最大承重(g)', '刀位允许长度(mm)']].to_dict()
    # print(machine_tool_magazine)

    machine_cool = machine_data[['支持冷却类型', '支持冷却液类型', '内冷却压力(bar)']].to_dict()
    # print(machine_cool)

    machine_precision = machine_data[
        ['定位精度(角分)', '重复定位精度(角分)', '反向差值', '回转定位精度(角分)', '回转重复定位精度(角分)']].to_dict()
    # print(machine_precision)

    # machine_power = {"主轴功率曲线": []}
    # machine_power_curve = machine_specifications["主轴功率曲线"].tolist()[1:]
    # point = [i for i in range(1, len(machine_power_curve) + 1)]
    # machine_power["主轴功率曲线"] = [point, machine_power_curve]
    # print(machine_power)}

    machine_power = {"主轴功率曲线": []}
    machine_power_curve = machine_data[
        ['主轴功率曲线1', '主轴功率曲线2', '主轴功率曲线3', '主轴功率曲线4', '主轴功率曲线5', '主轴功率曲线6',
         '主轴功率曲线7', '主轴功率曲线8', '主轴功率曲线9']].to_list()
    point = [i for i in range(1, 10)]
    machine_power["主轴功率曲线"] = [point, machine_power_curve]
    # print("主轴功率曲线", machine_power)

    # machine_torque = {"主轴扭矩曲线": []}
    # machine_torque_curve = machine_specifications["主轴扭矩曲线"].tolist()[1:]
    # point = [i for i in range(1, len(machine_torque_curve) + 1)]
    # machine_torque["主轴扭矩曲线"] = [point, machine_torque_curve]
    # # print(machine_torque)

    machine_torque = {"主轴扭矩曲线": []}
    machine_torque_curve = machine_data[
        ['主轴扭矩曲线1', '主轴扭矩曲线2', '主轴扭矩曲线3', '主轴扭矩曲线4', '主轴扭矩曲线5', '主轴扭矩曲线6',
         '主轴扭矩曲线7', '主轴扭矩曲线8', '主轴扭矩曲线9']].to_list()
    point = [i for i in range(1, 10)]
    machine_torque["主轴扭矩曲线"] = [point, machine_torque_curve]
    # print("主轴扭矩曲线", machine_torque)

    # machine_health = {"机床体检评分": []}
    # machine_health_curve = machine_specifications["机床体检评分"].tolist()[-9:]
    # # point = ["2022-06-01", "2022-06-02", "2022-06-03", "2022-06-04", "2022-06-05", "2022-06-06", "2022-06-07", "2022-06-08",
    # #          "2022-06-09"]
    # point = machine_specifications["机床体检日期"].tolist()[-9:]
    # machine_health["机床体检评分"] = [point, machine_health_curve]
    # # print(machine_health)

    temp_dict["task"] = str("update_machine_information")
    temp_dict["machine_id"] = str(machine_id)
    temp_dict["machine_basic_information"] = dict(machine_basic_information)  # 机器基本信息
    temp_dict["machine_constraints"] = dict(machine_constraints)  # 机器约束
    temp_dict["machine_tool_magazine"] = dict(machine_tool_magazine)  # 机器刀具库
    temp_dict["machine_cool"] = dict(machine_cool)  # 机器冷却
    temp_dict["machine_precision"] = dict(machine_precision)  # 机器精度
    temp_dict["machine_power"] = dict(machine_power)  # 机器功率
    temp_dict["machine_torque"] = dict(machine_torque)  # 机器扭矩
    # temp_dict["machine_health"] = dict(machine_health)  # 健康状态
    j = json.dumps(temp_dict, ensure_ascii=False)
    # print("生产资料： ", j)
    publish_message(STATIC_INFO, j)


def extract_machine_specifications(df):
    machine_basic_information = df.iloc[0, :6].to_dict()
    machine_constraints = df.iloc[0, 6:9].to_dict()
    machine_tool_magazine = df.iloc[0, 9:14].to_dict()
    machine_cool = df.iloc[0, 14:17].to_dict()
    machine_precision = df.iloc[0, 17:22].to_dict()

    # 主轴功率曲线、主轴扭矩曲线和机床体检评分是 Series 类型，需单独处理
    machine_power_curve = df.loc[0, "主轴功率曲线"]
    machine_torque_curve = df.loc[0, "主轴扭矩曲线"]
    machine_health_score = df.loc[0, "机床体检评分"]

    return (
        machine_basic_information,
        machine_constraints,
        machine_tool_magazine,
        machine_cool,
        machine_precision,
        machine_power_curve,
        machine_torque_curve,
        machine_health_score,
    )


def to_csv():
    path = 'conf/'
    # # 写入mock_machines.csv
    # df_machine_list = pd.DataFrame({
    #     'machine_list': machine_list,
    # })
    # df_machine_list.to_csv(path + 'mock_machine_list.csv', index=False)
    #
    # df_machine_type = pd.DataFrame({'machine_type': machine_type_list})
    # df_machine_type.to_csv(path + 'mock_machine_type.csv', index=False)
    #
    # df_machine_status = pd.DataFrame({'machine_status': machine_status})
    # df_machine_status.to_csv(path + 'mock_machine_status.csv', index=False)
    #
    # # 写入mock_tools.csv
    # df_tools = pd.DataFrame({'tool_id': tool_list})
    # df_tools.to_csv(path + 'mock_tool_list.csv', index=False)
    #
    # # 写入mock_programs.csv
    # df_program = pd.DataFrame({'program_name': program_list})
    # df_program.to_csv(path + 'mock_program_list.csv', index=False)
    #
    # # 写入mock_time_intervals.csv
    # df_time_intervals = pd.DataFrame({'time_interval': time_interval})
    # df_time_intervals.to_csv(path + 'mock_time_interval.csv', index=False)
    #
    # 写入machine_specifications.csv
    # 处理特殊的字典，将二维列表转换为 Series
    # machine_power_values = machine_power.pop("主轴功率曲线")
    # machine_power_series = pd.Series(data=machine_power_values[1], index=machine_power_values[0], name="主轴功率曲线")
    #
    # machine_torque_values = machine_torque.pop("主轴扭矩曲线")
    # machine_torque_series = pd.Series(data=machine_torque_values[1], index=machine_torque_values[0],
    #                                   name="主轴扭矩曲线")
    #
    # machine_health_values = machine_health.pop("机床体检评分")
    # machine_health_series = pd.Series(data=machine_health_values[1], index=machine_health_values[0],
    #                                   name="机床体检评分")
    #
    # # 创建单行 DataFrame
    # df_machine_basic_information = pd.DataFrame(machine_basic_information, index=[0])
    # df_machine_constraints = pd.DataFrame(machine_constraints, index=[0])
    # df_machine_tool_magazine = pd.DataFrame(machine_tool_magazine, index=[0])
    # df_machine_cool = pd.DataFrame(machine_cool, index=[0])
    # df_machine_precision = pd.DataFrame(machine_precision, index=[0])
    #
    # # 将 Series 添加为单列 DataFrame
    # df_machine_power = pd.DataFrame({"主轴功率曲线": machine_power_series})
    # df_machine_torque = pd.DataFrame({"主轴扭矩曲线": machine_torque_series})
    # df_machine_health = pd.DataFrame({"机床体检评分": machine_health_series})
    #
    # df_machine_specifications = pd.concat([
    #     df_machine_basic_information,
    #     df_machine_constraints,
    #     df_machine_tool_magazine,
    #     df_machine_cool,
    #     df_machine_precision,
    #     df_machine_power,
    #     df_machine_torque,
    #     df_machine_health,
    # ], axis=1)
    #
    # df_machine_specifications.to_csv(path + 'mock_machine_specifications.csv', index=False)
    #
    # # 写入pie_data.csv
    # df_pie_list = pd.DataFrame(pie_list)
    # df_pie_list.to_csv(path + 'mock_pie_data.csv', index=False)
    #
    # # 写入finish_rate_data.csv
    # df_finish_rate_list = pd.DataFrame(finish_rate_list)
    # df_finish_rate_list.to_csv(path + 'mock_finish_rate_list.csv', index=False)
    #
    # # 写入list_seven_data.csv
    # df_list_seven = pd.DataFrame(list_seven)
    # df_list_seven.to_csv(path + 'mock_list_seven.csv', index=False)
    #
    # # 写入product_id_list.csv
    # df_product_id_list = pd.DataFrame({'product_id': product_id_list})
    # df_product_id_list.to_csv(path + 'mock_product.csv', index=False)
    #
    # # 写入tool_info.csv
    # df_tool_info = pd.DataFrame(tool_info).T
    # df_tool_info.to_csv(path + 'mock_tool_info.csv', index=False)
    #
    # # 写入workpiece_info.csv
    # df_workpiece_info = pd.DataFrame(workpiece).T
    # df_workpiece_info.to_csv(path + 'mock_workpiece.csv', index=False)
    # status_mapping = pd.read_csv(os.path.join('conf', 'mock_machine_status_mapping.csv'))
    # MACHINE_STATUS = status_mapping.iloc[0].to_dict()
    # print(MACHINE_STATUS)
    # MACHINE_STATUS_INV = {v: k for k, v in MACHINE_STATUS.items()}
    # print(MACHINE_STATUS_INV)
    # 给定的数据
    # data = [
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
    #
    # flat_data = []
    # for entry in data:
    #     for day_4hour, week_1day, month_5day in zip(entry['day-4hour'], entry['week-1day'], entry['month-5day']):
    #         for i in range(len(day_4hour)):
    #             row = {'day-4hour': day_4hour[i],
    #                    'week-1day': week_1day[i],
    #                    'month-5day': month_5day[i]}
    #             flat_data.append(row)
    #
    # # 创建DataFrame
    # df = pd.DataFrame(flat_data)
    #
    # # 将DataFrame写入CSV文件
    # df.to_csv('mock_Electricity_consumption.csv', index=False)
