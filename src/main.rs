use std::str::FromStr;

struct EnvOpts {
    dir_unit: Option<String>,
    dir_log: Option<String>,
    quick_exit: bool,
}

fn get_env_opts() -> EnvOpts {
    let mut opts = EnvOpts {
        dir_unit: Some(String::from("/etc/minit.d")),
        dir_log: None,
        quick_exit: false,
    };

    match std::env::var("MINIT_UNIT_DIR") {
        Ok(val) => {
            if !val.eq_ignore_ascii_case("none") {
                opts.dir_unit = Some(val)
            }
        }
        _ => {}
    }

    match std::env::var("MINIT_LOG_DIR") {
        Ok(val) => {
            if !val.eq_ignore_ascii_case("none") {
                opts.dir_log = Some(val)
            }
        }
        _ => {}
    }

    match std::env::var("MINIT_QUICK_EXIT") {
        Ok(val) => match <bool as FromStr>::from_str(&val) {
            Ok(v) => opts.quick_exit = v,
            _ => {}
        },
        _ => {}
    }

    opts
}

fn main() {
    let _opts = get_env_opts();
    println!("minit: starting (#{})", env!("MINIT_COMMIT"));
}
