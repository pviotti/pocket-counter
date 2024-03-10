// Sample data: Replace this with your actual data
const data = [
    { date: "2023-08-01", count: 5 },
    { date: "2023-08-02", count: 10 },
    // ... more data entries
];

// Extract dates and counts for Chart.js
const dates = data.map(entry => entry.date);
const counts = data.map(entry => entry.count);

// Create a chart using Chart.js
const ctx = document.getElementById('chartCanvas');
new Chart(ctx, {
    type: 'line',
    data: {
        labels: dates,
        datasets: [{
            label: 'Unread Articles Count',
            data: counts,
            borderColor: 'blue',
            fill: false,
        }]
    },
    // options: {
    //     scales: {
    //         x: {
    //             type: 'time',
    //             time: {
    //                 unit: 'day',
    //                 // displayFormats: {
    //                 //     day: 'MMM D',
    //                 // },
    //             },
    //             title: {
    //                 display: true,
    //                 text: 'Date'
    //             }
    //         },
    //         y: {
    //             beginAtZero: true,
    //             title: {
    //                 display: true,
    //                 text: 'Unread Articles Count'
    //             }
    //         }
    //     }
    // }
});
