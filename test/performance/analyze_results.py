import json
import matplotlib.pyplot as plt
import pandas as pd
from datetime import datetime
import sys
import re

def clean_value(value):
    """Clean and convert k6 metric values"""
    # Remove units and convert to float
    value = value.strip()
    if 'µs' in value:
        # Convert microseconds to milliseconds
        return float(value.replace('µs', '')) / 1000
    elif 'ms' in value:
        # Already in milliseconds
        return float(value.replace('ms', ''))
    elif 's' in value:
        # Convert seconds to milliseconds
        return float(value.replace('s', '')) * 1000
    elif '%' in value:
        # Convert percentage to decimal
        return float(value.replace('%', '')) / 100
    elif '/' in value:
        # Handle rates (e.g., "4960.674276/s")
        return float(value.split('/')[0])
    try:
        return float(value)
    except ValueError:
        return value

def parse_k6_metrics(filename):
    with open(filename, 'r') as f:
        content = f.read()

    metrics = {}

    # Parse http_req_duration
    duration_match = re.search(r'http_req_duration.*?avg=([\d.]+µs)\s+min=([\d.]+µs)\s+med=([\d.]+µs)\s+max=([\d.]+ms)\s+p\(90\)=([\d.]+µs)\s+p\(95\)=([\d.]+µs)', content)
    if duration_match:
        metrics['http_req_duration'] = {
            'avg': clean_value(duration_match.group(1)),
            'min': clean_value(duration_match.group(2)),
            'med': clean_value(duration_match.group(3)),
            'max': clean_value(duration_match.group(4)),
            'p90': clean_value(duration_match.group(5)),
            'p95': clean_value(duration_match.group(6))
        }

    # Parse iterations
    iterations_match = re.search(r'iterations\.*?:\s*([\d,]+)\s', content)
    if iterations_match:
        count = iterations_match.group(1).replace(',', '')
        metrics['iterations'] = {
            'count': int(float(count))
        }

    # Parse error rate
    http_errors_match = re.search(r'http_req_failed.*?:\s*([\d.]+)%', content)
    if http_errors_match:
        metrics['errors'] = {
            'rate': float(http_errors_match.group(1)) / 100
        }
    else:
        # Default to 0 if no errors found
        metrics['errors'] = {'rate': 0.0}

    return metrics

def analyze_test_results(filename):
    # Parse k6 output
    metrics = parse_k6_metrics(filename)
    
    # Basic statistics
    stats = {
        'Total Requests': metrics['iterations']['count'],
        'Error Rate': f"{metrics['errors']['rate']*100:.2f}%",
        'Avg Response Time': f"{metrics['http_req_duration']['avg']:.2f}ms",
        '95th Percentile': f"{metrics['http_req_duration']['p95']:.2f}ms",
        'Max Response Time': f"{metrics['http_req_duration']['max']:.2f}ms",
    }

    # Create summary report
    report = [
        "Performance Test Results",
        "=====================\n",
        f"Test Date: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n",
        "Key Metrics:",
        f"- Total Requests: {stats['Total Requests']}",
        f"- Error Rate: {stats['Error Rate']}",
        f"- Average Response Time: {stats['Avg Response Time']}",
        f"- 95th Percentile: {stats['95th Percentile']}",
        f"- Maximum Response Time: {stats['Max Response Time']}\n",
        "Performance Requirements Analysis:",
        f"- Response Time Requirement (< 3ms): {'PASSED' if metrics['http_req_duration']['p95'] < 3 else 'FAILED'}",
        f"- Error Rate Requirement (< 1%): {'PASSED' if metrics['errors']['rate'] < 0.01 else 'FAILED'}",
        f"- Throughput Requirement (5000 req/s): {'PASSED' if metrics['iterations']['count']/30 >= 4950 else 'FAILED'}\n",
        "Detailed Analysis:",
        "- Response Time Distribution:",
        f"  * 50th percentile (median): {metrics['http_req_duration']['med']:.2f}ms",
        f"  * 90th percentile: {metrics['http_req_duration']['p90']:.2f}ms",
        f"  * 95th percentile: {metrics['http_req_duration']['p95']:.2f}ms",
        f"  * Average: {metrics['http_req_duration']['avg']:.2f}ms"
    ]

    # Generate plots
    plt.figure(figsize=(12, 6))
    percentiles = ['Median', 'p90', 'p95', 'Average']
    values = [
        metrics['http_req_duration']['med'],
        metrics['http_req_duration']['p90'],
        metrics['http_req_duration']['p95'],
        metrics['http_req_duration']['avg']
    ]
    
    bars = plt.bar(percentiles, values, color='skyblue')
    plt.axhline(y=3, color='r', linestyle='--', label='3ms Requirement')
    
    # Add value labels on top of each bar
    for bar in bars:
        height = bar.get_height()
        plt.text(bar.get_x() + bar.get_width()/2., height,
                f'{height:.3f}ms',
                ha='center', va='bottom')
    
    plt.title('Response Time Distribution')
    plt.ylabel('Time (ms)')
    plt.legend()
    plt.grid(True, axis='y', linestyle='--', alpha=0.7)
    plt.savefig('response_time_percentiles.png', bbox_inches='tight', dpi=300)
    plt.close()

    # Write report
    with open('performance_report.txt', 'w') as f:
        f.write('\n'.join(report))

    return stats

if __name__ == '__main__':
    if len(sys.argv) > 1:
        filename = sys.argv[1]
    else:
        filename = 'k6_output.txt'
    
    try:
        stats = analyze_test_results(filename)
        print("\n=== Performance Analysis Results ===")
        print(f"\nKey Performance Metrics:")
        print(f"- Total Requests: {stats['Total Requests']}")
        print(f"- Error Rate: {stats['Error Rate']}")
        print(f"- Average Response Time: {stats['Avg Response Time']}")
        print(f"- 95th Percentile: {stats['95th Percentile']}")
        print(f"\nDetailed results have been saved to:")
        print("- performance_report.txt")
        print("- response_time_percentiles.png")
    except Exception as e:
        print(f"Error during analysis: {e}")
        import traceback
        print(traceback.format_exc())
