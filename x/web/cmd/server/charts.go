package main

var (
	chartsTemplate = `
{{define "charts"}}
<div class="row">
	<div class="col-lg-6">
		<div class="panel panel-default">
			{{template "boxes" .Chart.Box}}
		</div>
	</div>
	<div class="col-lg-6">
		<div class="panel panel-default">
			{{template "scatter" .Chart.Scatter}}
		</div>
	</div>
</div>	
{{end}}
	`

	boxesTemplate = `
{{define "boxes"}}
			<script>var d = {{.}} </script>
			<div id="boxplot" />
				<script>
				Highcharts.chart('boxplot', {
				
				       chart: {
				           type: 'boxplot'
				       },
				
				       title: {
				           text: '{{.Title}}'
				       },
				
				       legend: {
				           enabled: false
				       },
				
				       xAxis: {
				           categories: {{.XLabels}},
				           title: {
				               text: '{{.XTitle}}'
				           }
				       },
				
				       yAxis: {
				           title: {
				               text: 'Fitness'
				           },
				           min: {{.YMin}},
				           max: {{.YMax}}
				       },
				
				       series: [{
				           name: 'Fitness',
				           data: {{.Data}},
				           tooltip: {
				               headerFormat: '<em>Trial {point.key}</em><br/>'
				           }
				       }]
				
				   });
				   
				</script>
			</div>
{{end}}
	`

	scatterTemplate = `
{{define "scatter"}}
<script>var d = {{.}} </script>
			<div id="scatter" />
				<script>
     Highcharts.chart('scatter', {
        chart: {
            type: 'scatter',
            zoomType: 'xy'
        },
        title: {
            text: 'Fitness by Complexity'
        },
        xAxis: {
            title: {
                enabled: true,
                text: 'Complexity'
            },
            startOnTick: true,
            endOnTick: true,
            showLastLabel: true
        },
        yAxis: {
            title: {
                text: 'Fitness'
            }
        },
        plotOptions: {
            scatter: {
                marker: {
                    radius: 5,
                    states: {
                        hover: {
                            enabled: true,
                            lineColor: 'rgb(100,100,100)'
                        }
                    }
                },
                states: {
                    hover: {
                        marker: {
                            enabled: false
                        }
                    }
                },
                tooltip: {
                    headerFormat: '<b>{series.name}</b><br>',
                    pointFormat: 'Complexity: {point.x}, Fitness: {point.y}'
                }
            }
        },
        series: [{
            name: 'Best',
            color: 'rgba(223, 83, 83, .5)',
            data: {{.Data}}

        }]
    });
    			</script>

			</div>
{{end}}
`
)
