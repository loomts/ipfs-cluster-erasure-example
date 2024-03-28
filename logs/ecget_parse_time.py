import re
from datetime import datetime, timezone, timedelta

# 打开日志文件
with open('newfile', 'r') as file:
    lines = file.readlines()

# 存储时间差结果
time_diffs = []

# 奇数行和偶数行配对
for odd_line, even_line in zip(lines[::2], lines[1::2]):
    odd_times = re.findall(r'\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{4}', odd_line)
    even_times = re.findall(r'\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}', even_line)

    # 设置时区为UTC+8
    tz = timezone(timedelta(hours=8))

    # 计算时间差
    results = []
    for odd_time, even_time in zip(odd_times, even_times):
        odd_dt = datetime.strptime(odd_time, '%Y-%m-%dT%H:%M:%S.%f%z')
        # 把偶数行的时间转换为带时区的时间
        even_dt = datetime.strptime(even_time, '%Y/%m/%d %H:%M:%S').replace(tzinfo=tz)
        even_dt += (timedelta(seconds=1))
        time_diff = even_dt - odd_dt
        results.append(time_diff.total_seconds())

    print(','.join(map(str, results)), end="\n")