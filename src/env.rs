use std::env;
use std::str::FromStr;

pub fn env_str(key: &str, out: &mut Option<String>) {
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

pub fn env_bool(key: &str, out: &mut bool) {
    match env::var(key) {
        Ok(val) => match <bool as FromStr>::from_str(val.as_str()) {
            Ok(val) => *out = val,
            _ => {}
        },
        _ => {}
    }
}
