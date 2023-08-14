use std::env;
use std::str::FromStr;

fn env_str(key: &str, out: &mut Option<String>) {
    match env::var(key) {
        Ok(val) => {
            if val.eq_ignore_ascii_case("none") {
                *out = None
            } else {
                *out = Some(val);
            }
        }
        _ => {}
    }
}

fn env_bool(key: &str, out: &mut bool) {
    match env::var(key) {
        Ok(val) => match <bool as FromStr>::from_str(val.as_str()) {
            Ok(val) => *out = val,
            _ => {}
        },
        _ => {}
    }
}

fn main() {
    let mut opt_dir_unit: Option<String> = Some(String::from("/etc/minit.d"));
    let mut opt_dir_log: Option<String> = None;
    let mut opt_quick_exit = false;

    env_str("MINIT_UNIT_DIR", &mut opt_dir_unit);
    env_str("MINIT_LOG_DIR", &mut opt_dir_log);
    env_bool("MINIT_QUICK_EXIT", &mut opt_quick_exit);

    println!("minit: starting (#{})", env!("MINIT_COMMIT"));
}
