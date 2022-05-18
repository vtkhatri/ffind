use std::env;
use std::io;
use std::process;

#[derive(Debug, PartialEq, Eq)]
enum DebugLevel {
    DEBUG,
    NONE,
}

struct DebugManager {
    debug_level: DebugLevel,
}

trait DebugPrints {
    fn print(&self, to_print: String);
}

impl DebugPrints for DebugManager {
    fn print(&self, to_print: String) {
        if self.debug_level == DebugLevel::DEBUG {
            eprintln!("{}", to_print);
        }
        return;
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();

    // making the command
    let cmd_args_result = make_command(args);
    let cmd_args = match cmd_args_result {
        Err(e) => {
            println!("Error making command: {}", e);
            process::exit(1);
        }
        Ok(cmd_args) => cmd_args,
    };

    // executing the command
    let exec_status_result = process::Command::new("find")
        .args(cmd_args)
        .stdout(process::Stdio::inherit())
        .stderr(process::Stdio::inherit())
        .status();

    // proper error code propogation
    // spaghettified because Command::status() returns - Ok(Exitstatus(Exitstatus(code)))
    // we need to unwrap highest Ok then unwrap the Exitstatus to get to the return code of
    // command executed
    process::exit(match exec_status_result {
        Ok(exec_out) => match exec_out.code() {
            Some(code) => code,
            None => {
                println!("Process terminated by signal");
                130 // because when I press ctrl-c and echo $? right after, it prints 130
            }
        },
        Err(e) => {
            println!("Error executing command: {}", e);
            2
        }
    });
}

fn print_usage() {
    println!("Usage: ffind [-fdrih] [-e=maxdepth] [--debug --help] [expression] [path]");
}

#[derive(Debug)]
struct SortedArgs {
    short_args: Vec<String>,
    long_args: String,
    file_name: String,
    path: Vec<String>,
    exec_args: Vec<String>,
}

fn make_command(args_in: Vec<String>) -> Result<Vec<String>, io::Error> {
    let mut args_out = Vec::new();
    let sorted_args = sort_args(args_in)?;

    args_out.append(&mut sorted_args.path.to_vec());
    args_out.append(&mut sorted_args.short_args.to_vec());
    args_out.push(sorted_args.file_name);
    args_out.append(&mut sorted_args.exec_args.to_vec());

    Ok(args_out)
}

fn sort_args(args_in: Vec<String>) -> Result<SortedArgs, io::Error> {
    let mut args_in_for_exec = args_in.to_vec();
    let mut short_args: Vec<String> = Vec::new();
    let mut long_args: Vec<String> = Vec::new();
    let mut file_name = String::from("");
    let mut path = Vec::new();
    let mut exec_args = Vec::new();

    for (arg_no, arg) in args_in.iter().enumerate() {
        match arg.rsplit_once('-') {
            // separate out --flags, -flags, and -exec among arguments
            Some((first, second)) => {
                if first == "-" {
                    long_args.push(second.to_string());
                } else {
                    if second == "exec" {
                        exec_args = args_in_for_exec.drain(arg_no..).collect::<Vec<String>>();
                        break;
                    } else {
                        short_args.push(second.to_string());
                    }
                }
            }
            // only 2 things don't begin with '-' => filename and path, filename always preceding
            None => {
                // ffind itself is argument no. 0, so we avoid checking that
                if arg_no != 0 {
                    if file_name.is_empty() {
                        file_name = arg.to_string();
                    } else {
                        path.push(arg.to_string());
                    }
                }
            }
        }
    }

    let debug_man = &mut DebugManager {
        debug_level: DebugLevel::NONE,
    };

    let sorted_args = SortedArgs {
        long_args: process_long_args(long_args, debug_man)?,
        short_args: get_short_args(short_args)?,
        file_name: file_name,
        path: path,
        exec_args: exec_args,
    };

    debug_man.print(format!("args_in={:?}", args_in));
    debug_man.print(format!("sorted_args={:?}", sorted_args));

    Ok(sorted_args)
}

fn get_short_args(args_in: Vec<String>) -> Result<Vec<String>, io::Error> {
    let mut ret_short_args: Vec<String> = Vec::new();
    enum GlobType {
        CaseSensitiveName,
        CaseSensitiveRegex,
        CaseInsensitiveName,
        CaseInsensitiveRegex,
    }
    let mut glob_type = GlobType::CaseSensitiveName;

    for args in args_in {
        for arg in args.chars() {
            match arg {
                'f' => {
                    ret_short_args.push(String::from("-type"));
                    ret_short_args.push(String::from("f"));
                }
                'd' => {
                    ret_short_args.push(String::from("-type"));
                    ret_short_args.push(String::from("d"));
                }
                'e' => {
                    match args.rsplit_once('=') {
                        Some((_, depth)) => match depth.parse::<i32>() {
                            Ok(_) => {
                                ret_short_args.insert(0, String::from("-maxdepth"));
                                ret_short_args.insert(1, String::from(depth));
                            }
                            Err(e) => {
                                return Err(io::Error::new(
                                    io::ErrorKind::Other,
                                    format!("-e flag needs integer maxdepth: {}", e),
                                ));
                            }
                        },
                        None => {
                            return Err(io::Error::new(
                                io::ErrorKind::Other,
                                "-e flag used without maxdepth",
                            ));
                        }
                    }
                    break;
                }
                'r' => {
                    glob_type = match glob_type {
                        GlobType::CaseSensitiveName => GlobType::CaseSensitiveRegex,
                        GlobType::CaseInsensitiveName => GlobType::CaseInsensitiveRegex,
                        _ => glob_type,
                    }
                }
                'i' => {
                    glob_type = match glob_type {
                        GlobType::CaseSensitiveName => GlobType::CaseInsensitiveName,
                        GlobType::CaseSensitiveRegex => GlobType::CaseInsensitiveRegex,
                        _ => glob_type,
                    }
                }
                'h' => {
                    print_usage();
                    process::exit(0);
                }
                _ => {
                    return Err(io::Error::new(
                        io::ErrorKind::Other,
                        format!("flag -{} not recognized", arg),
                    ));
                }
            }
        }
    }

    // If we add filename globtype here, the filename can be added as-is in make_command call
    match glob_type {
        GlobType::CaseSensitiveName => ret_short_args.push(String::from("-name")),
        GlobType::CaseSensitiveRegex => ret_short_args.push(String::from("-regex")),
        GlobType::CaseInsensitiveName => ret_short_args.push(String::from("-iname")),
        GlobType::CaseInsensitiveRegex => ret_short_args.push(String::from("-iregex")),
    }

    return Ok(ret_short_args);
}

fn process_long_args(
    arg_in: Vec<String>,
    mut debug_man: &mut DebugManager,
) -> Result<String, io::Error> {
    for arg in arg_in {
        match arg.as_str() {
            "help" => {
                print_usage();
                process::exit(0);
            }
            "debug" => {
                debug_man.debug_level = DebugLevel::DEBUG;
            }
            _ => {
                return Err(io::Error::new(
                    io::ErrorKind::Other,
                    format!("flag --{} not recognized", arg),
                ));
            }
        }
    }
    Ok(String::from("long arguments are processed"))
}
