import * as React from "react";
import {useState} from "react";
import ConfigObject from "./ConfigObject";

export default function ConfigSensor(props) {
    const [conf, setConf] = React.useState({
        addr: "127.0.0.1",
        period: "1200"
    })

    const [settings, setSettings] = React.useState({})


    const [confMeta, setConfMeta] = React.useState([
        {
            name: "addr",
            display_name: "Address",
            type: "text"
        },
        {
            name: "period",
            display_name: "Period, ms",
            type: "number",
            min_value: "1",
            max_value: "10000"
        },
        {
            name: "tags",
            display_name: "Modbus tags",
            type: "table",
            columns: [
                {
                    name: "tag_name",
                    display_name: "Name",
                    type: "text",
                    defaultValue: "tagName"
                },
                {
                    name: "tag_addr",
                    display_name: "Address",
                    type: "number",
                    min_value: "1",
                    max_value: "254",
                    defaultValue: 42
                },
                {
                    name: "tags_inner",
                    display_name: "Modbus inner tags",
                    type: "table",
                    columns: [
                        {
                            name: "tag_name1",
                            display_name: "Name1",
                            type: "text"
                        },
                        {
                            name: "tag_addr1",
                            display_name: "Address1",
                            type: "number",
                            min_value: "1",
                            max_value: "254"
                        }
                    ]
                }

            ]
        }
    ])

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
            <div>{JSON.stringify(conf)}</div>
        </div>
    )
}
