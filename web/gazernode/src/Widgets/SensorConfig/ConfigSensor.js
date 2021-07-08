import * as React from "react";
import {useState} from "react";
import ConfigObject from "./ConfigObject";

export default function ConfigSensor(props) {
    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        setFirstRendering(false)
    }

    const drawItem = (meta, data) => {
        return (
            <ConfigObject Meta={meta} Data={data} OnChangedValue={(n,v) => {
                props.OnConfigChanged(v)
            }}/>
        )
    }

    return (
        <div>
            <div>{drawItem(props.ConfigMeta, props.Config)}</div>
        </div>
    )
}
