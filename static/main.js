function stringToColor(str) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    let color = '#';
    for (let i = 0; i < 3; i++) {
        const value = (hash >> (i * 8)) & 0xFF;
        color += ('00' + value.toString(16)).substr(-2);
    }
    return color;
}

document.addEventListener('DOMContentLoaded', function() {
    const diskTableBody = document.getElementById('disk-table-body');
    const ctx = document.getElementById('temperatureChart').getContext('2d');
    let temperatureChart = new Chart(ctx, {
        type: 'line',
        data: {
            datasets: []
        },
        options: {
            scales: {
                x: {
                    type: 'time',
                    time: {
                        unit: 'second',
                        displayFormats: {
                            second: 'HH:mm:ss'
                        }
                    },
                    title: {
                        display: true,
                        text: 'Time'
                    }
                },
                y: {
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: 'Temperature (°C)'
                    }
                }
            }
        }
    });

    function fetchAndUpdateDisks() {
        fetch('/v1/api/disks?temp=true')
            .then(response => response.json())
            .then(data => {
                updateDiskTable(data.disks);
            })
            .catch(error => console.error('Error fetching disk data:', error));
    }

    function updateDiskTable(disks) {
        let tableHTML = '';
        disks.forEach(disk => {
            tableHTML += `
                <tr>
                    <td>${disk.Name}</td>
                    <td>${disk.Transport}</td>
                    <td>${disk.Size}</td>
                    <td>${disk.Model}</td>
                    <td>${disk.Serial}</td>
                    <td>${disk.Type}</td>
                    <td>${disk.Temperature}</td>
                </tr>
            `;
        });
        diskTableBody.innerHTML = tableHTML;
    }

    function fetchAndUpdateTemperatureChart() {
        fetch('/v1/api/snapshots')
            .then(response => response.json())
            .then(snapshots => {
                updateTemperatureChart(snapshots);
            })
            .catch(error => console.error('Error fetching temperature data:', error));
    }

    function updateTemperatureChart(snapshots) {
        // Clear existing datasets
        temperatureChart.data.datasets = [];

        snapshots.forEach(snapshot => {
            const time = new Date(snapshot.TimeStamp);
            snapshot.HDD.forEach(disk => {
                let dataset = temperatureChart.data.datasets.find(d => d.label === disk.Name);
                if (!dataset) {
                    dataset = {
                        label: disk.Name,
                        data: [],
                        fill: false,
                        borderColor: stringToColor(disk.Name),
                        borderWidth: 1
                    };
                    temperatureChart.data.datasets.push(dataset);
                }

                dataset.data.push({
                    x: time,
                    y: disk.Temperature
                });
            });
        });

        temperatureChart.update();
    }

    // Chart.js zoom and pan configuration
    temperatureChart.options.plugins.zoom = {
        zoom: {
            wheel: {
                enabled: true,
            },
            pinch: {
                enabled: true
            },
            mode: 'x',
        },
        pan: {
            enabled: true,
            mode: 'x',
        }
    };

    fetchAndUpdateDisks();
    fetchAndUpdateTemperatureChart();
    setInterval(fetchAndUpdateDisks, 5000);
    setInterval(fetchAndUpdateTemperatureChart, 5000);
});
