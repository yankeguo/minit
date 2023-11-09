use crate::types::GeneralResult;
use crate::utils::str_to_bool;
use serde::{Deserialize, Serialize};
use std::collections::{BTreeMap, HashMap};
use std::fs;
use std::path::Path;

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

impl Unit {
    pub fn new() -> Self {
        Unit {
            kind: UnitKind::Once,
            name: "".to_string(),
            group: None,
            count: None,
            dir: None,
            shell: None,
            env: None,
            command: None,
            charset: None,
            raw: None,
            files: None,
            cron: None,
            immediate: None,
        }
    }
}

pub fn load_arg_unit(
    result: &mut Vec<Unit>,
    args: &mut dyn Iterator<Item = String>,
) -> GeneralResult<()> {
    let mut args: Vec<String> = args.skip(1).collect();
    let mut opts: Vec<String> = vec![];
    for (i, arg) in args.iter().enumerate() {
        if arg == "--" {
            opts = args[0..i].iter().cloned().collect();
            args = args[i + 1..].iter().cloned().collect();
            break;
        }
    }
    let arg0 = Path::new(&args[0]);
    if let Some(file_name) = arg0.file_name() {
        if file_name.eq_ignore_ascii_case("minit") {
            args = args[1..].iter().cloned().collect();
        }
    }
    if args.is_empty() {
        return Ok(());
    }
    let mut unit = Unit::new();
    unit.kind = UnitKind::Daemon;
    unit.name = String::from("arg-main");
    unit.command = Some(args);
    for opt in opts {
        if opt.ends_with("-once") {
            unit.kind = UnitKind::Once;
        }
    }
    result.push(unit);
    return Ok(());
}

pub fn load_env_unit(
    result: &mut Vec<Unit>,
    envs: &mut dyn Iterator<Item = (String, String)>,
) -> GeneralResult<()> {
    let envs: HashMap<String, String> = HashMap::from_iter(envs);

    let mut unit = Unit::new();
    unit.kind = UnitKind::Daemon;
    unit.name = envs
        .get("MINIT_MAIN_NAME")
        .cloned()
        .unwrap_or("env-main".to_string());
    unit.group = envs.get("MINIT_MAIN_GROUP").cloned();
    unit.dir = envs.get("MINIT_MAIN_DIR").cloned();
    unit.charset = envs.get("MINIT_MAIN_CHARSET").cloned();

    if let Some(env_main) = envs.get("MINIT_MAIN") {
        //TODO: split command
        unit.command = Some(vec![env_main.clone()]);
    }

    if let Some(kind) = envs.get("MINIT_MAIN_KIND") {
        if kind.eq_ignore_ascii_case("daemon") {
            unit.kind = UnitKind::Daemon;
        } else if kind.eq_ignore_ascii_case("once") {
            unit.kind = UnitKind::Once;
        } else if kind.eq_ignore_ascii_case("cron") {
            unit.kind = UnitKind::Cron;
            if let Some(cron) = envs.get("MINIT_MAIN_CRON") {
                unit.cron = Some(cron.clone());
            } else {
                return Err("cron unit must have cron expression".into());
            }
            if let Some(immediate) = envs.get("MINIT_MAIN_IMMEDIATE") {
                unit.immediate = Some(str_to_bool(immediate));
            }
        } else {
            return Err(format!("unknown kind: {}", kind).into());
        }
    } else if let Some(once) = envs.get("MINIT_MAIN_ONCE") {
        if str_to_bool(once) {
            unit.kind = UnitKind::Once;
        }
    }

    if unit.command.is_some() {
        result.push(unit);
    }
    return Ok(());
}

pub fn load_dir_units(result: &mut Vec<Unit>, dir: &String) -> GeneralResult<()> {
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
    return Ok(());
}
