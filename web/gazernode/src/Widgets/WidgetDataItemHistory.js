import React, {useEffect, useState} from 'react';
import Chart from "chart.js";
import Grid from "@material-ui/core/Grid";
import ButtonGroup from "@material-ui/core/ButtonGroup";
import Button from "@material-ui/core/Button";
import Request from "../request";

// LineChart
class LineChart extends React.Component {
    constructor(props) {
        super(props);
        this.canvasRef = React.createRef();
    }

    componentDidUpdate() {
        this.myChart.data.labels = this.props.data.map(d => d.time);
        this.myChart.data.datasets[0].data = this.props.data.map(d => d.value);
        this.myChart.update();
    }

    componentDidMount() {
        const canvas = this.canvasRef.current;
        const ctx = canvas.getContext('2d');

        let gradient = ctx.createLinearGradient(0, 0, 100, 400);
        gradient.addColorStop(0, 'rgba(250,174,50,1)');
        gradient.addColorStop(1, 'rgba(250,174,50,0)');

        this.myChart = new Chart(this.canvasRef.current, {
            type: 'line',
            options: {
                maintainAspectRatio: false,
                scales: {
                    xAxes: [
                        {
                            type: 'time',
                            time: {
                                unit: 'minute'
                            }
                        }
                    ],
                    yAxes: [
                        {
                        }
                    ]
                },
                animation: {
                    duration: 0
                }
            },
            data: {
                labels: this.props.data.map(d => d.time),
                datasets: [{
                    label: this.props.title,
                    data: this.props.data.map(d => d.value),
                    fill: false,
                    fillColor : "#FF0000",
                    backgroundColor: this.props.color,
                    pointRadius: 0,
                    borderColor: this.props.color,
                    borderWidth: 1,
                    lineTension: 0
                }]
            }
        });
    }

    render() {
        return (
                <canvas ref={this.canvasRef} />
        );
    }
}

export default function WidgetDataItemHistory(props) {
    const [data, setData] = useState([])
    const [firstRendering, setFirstRendering] = useState(true)
    const [requestHistoryProcessing, setRequestHistoryProcessing] = useState(false)
    if (firstRendering) {
        setFirstRendering(false)
    }

    const isNumeric = (str) => {
        if (typeof str != "string") return false // we only process strings!
        return !isNaN(str) && // use type coercion to parse the _entirety_ of the string (`parseFloat` alone does not do this)...
            !isNaN(parseFloat(str)) // ...and ensure strings of whitespace fail
    }

    const requestHistory = (path) => {
        if (requestHistoryProcessing) {
            return
        }
        setRequestHistoryProcessing(true)

        let dtEnd = new Date().getTime() * 1000
        let dtBegin = dtEnd - 1000000 * 60

        if (rbTimeFilter === 0) {
            dtBegin = dtEnd - 1000000 * 60
        }
        if (rbTimeFilter === 1) {
            dtBegin = dtEnd - 1000000 * 60 * 5
        }
        if (rbTimeFilter === 2) {
            dtBegin = dtEnd - 1000000 * 60 * 15
        }
        if (rbTimeFilter === 3) {
            dtBegin = dtEnd - 1000000 * 60 * 30
        }
        if (rbTimeFilter === 4) {
            dtBegin = dtEnd - 1000000 * 60 * 60
        }

        let req = {
            "name": path,
            "dt_begin": dtBegin,
            "dt_end": dtEnd,
        }
        Request('data_item_history', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            let data = [];
                            for(let i = 0; i < result.history.items.length; i++) {
                                if (isNumeric(result.history.items[i].v) && result.history.items[i].u !== "error") {
                                    data.push({
                                        time: new Date(result.history.items[i].t / 1000),
                                        value: result.history.items[i].v
                                    });
                                }
                            }
                            setData(data)
                            setRequestHistoryProcessing(false)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
                setRequestHistoryProcessing(false)
            })
            .catch((err) => {
                setRequestHistoryProcessing(false)
                //setErrorMessage("Unknown error")
            })
    }


    useEffect(() => {
        const timer = setInterval(() => {
            requestHistory(props.Path)
        }, 500);
        return () => clearInterval(timer);
    });

    const [rbTimeFilter, setRbTimeFilter] = useState(0)

    const rbTimeFilterButtonVariant = (index) => {
        if (index === rbTimeFilter)
            return "contained"
        return "outlined"
    }

    //let ddd = getRandomDateArray(15)

    return (
        <Grid container direction="column">
            <Grid item>
                <ButtonGroup >
                    <Button variant={rbTimeFilterButtonVariant(0)} color="primary" onClick={() => {setRbTimeFilter(0)}}>1 min </Button>
                    <Button variant={rbTimeFilterButtonVariant(1)} color="primary" onClick={() => {setRbTimeFilter(1)}}>5 min </Button>
                    <Button variant={rbTimeFilterButtonVariant(2)} color="primary" onClick={() => {setRbTimeFilter(2)}}>15 min </Button>
                    <Button variant={rbTimeFilterButtonVariant(3)} color="primary" onClick={() => {setRbTimeFilter(3)}}>30 min</Button>
                    <Button variant={rbTimeFilterButtonVariant(4)} color="primary" onClick={() => {setRbTimeFilter(4)}}>60 min</Button>
                </ButtonGroup>
            </Grid>
            <Grid item>
                <LineChart title={"asd"} data={data} color="#52bdff"/>
            </Grid>
        </Grid>

    )
}
