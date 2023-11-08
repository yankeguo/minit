use serde::{Deserialize, Serialize};
use std::collections::BTreeMap;
use std::error::Error;
use std::fs;

#[derive(Serialize, Deserialize, PartialEq, Debug)]
pub enum UnitKind {
    #[serde(rename = "once")]
    Once,
    #[serde(rename = "daemon")]
    Daemon,
    #[serde(rename = "render")]
    Render,
    #[serde(rename = "cron")]
    Cron,
}

#[derive(Serialize, Deserialize, PartialEq, Debug)]
pub struct Unit {
    kind: UnitKind,
    name: String,
    group: Option<String>,
    count: Option<i32>,

    // execution options
    dir: Option<String>,
    shell: Option<String>,
    env: Option<BTreeMap<String, String>>,
    command: Option<Vec<String>>,
    charset: Option<String>,

    // render options
    raw: Option<bool>,
    files: Option<Vec<String>>,

    // cron options
    cron: Option<String>,
    immediate: Option<bool>,
}

pub fn load_units(dir: &String) -> Result<Vec<Unit>, Box<dyn Error>> {
    let mut result: Vec<Unit> = vec![];
    for entry in fs::read_dir(dir)? {
        let entry = entry?;
        let file_name = entry.file_name();
        let file_name = file_name.to_str().unwrap();
        if !file_name.ends_with(".yaml") && !file_name.ends_with(".yml") {
            continue;
        }
        let metadata = entry.metadata()?;
        if !metadata.is_file() {
            continue;
        }

        let content = fs::read_to_string(entry.path())?;
        for chunk in content.split("---") {
            let chunk = chunk.trim();
            if chunk.is_empty() {
                continue;
            }
            let unit: Unit = serde_yaml::from_str(chunk)?;
            result.push(unit)
        }
    }
    return Ok(result);
}
