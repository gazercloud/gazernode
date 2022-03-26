import React, {useEffect, useRef, useState} from 'react';
import NewTimeChart from "./TimeChart";


export default function WidgetTimeChart(props) {
    const element1 = useRef();

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        setFirstRendering(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
        }, 1000);
        return () => clearInterval(timer);
    });

    const render1 = () => {
        if (element1.current !== undefined) {
            //console.log("chart render", element1.current.id)
            let ch = NewTimeChart(element1.current)
            let a1 = ch.addArea("a1")
            let ser1 = a1.addSeries("ser1")
            ser1.data = props.Data
            ch.updateHeight()
            ch.setHorScale(props.MinTime, props.MaxTime)
            ch.draw()

        }
    }

    render1()

    return (
        <div style={{display: "block"}}>
            <canvas style={{backgroundColor: "#111"}} id={"qwe"} ref={element1} width={props.ChartWidth}/>
        </div>
    );

}

