import re

# Initialize the total values
total_recon_time_diff = 0
total_recon_size_diff = 0

# Open the file and read line by line
with open('parse-addec-ecget-same.log', 'r') as file:
    i = 0
    for line in file:
        # Use regex to find the values
        match = re.search(r'rs_recon_time_diff:(\d+\.\d+)ms, rs_recon_size_diff:(\d+)', line)
        if match:
            i += 1
            total_recon_time_diff += float(match.group(1))
            total_recon_size_diff += int(match.group(2))
            if i % 2 == 0:
                # Calculate re_recon_rate_diff
                if total_recon_time_diff != 0:
                    re_recon_rate_diff = total_recon_size_diff / (total_recon_time_diff / 1000)  # convert ms to s
                    print(f"rs_recon_time_diff: {total_recon_time_diff},",
                          f"rs_recon_size_diff:{total_recon_size_diff},",
                          f"re_recon_rate_diff: {re_recon_rate_diff}")
                total_recon_time_diff = 0
                total_recon_size_diff = 0
