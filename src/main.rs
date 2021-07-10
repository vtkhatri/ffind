use std::env;
use std::process;
use std::io;

fn main() {
    let args: Vec<String> = env::args().collect();

    // making the command
    let cmd_args_result = make_command(args);
    let cmd_args = match cmd_args_result {
        Err(e) => {
            println!("Error making command: {}", e);
            process::exit(1);
        },
        Ok(cmd_args) => cmd_args,
    };
/*    
    // executing the command
    let exec_status_result = process::Command::new("find")
                                              .args(cmd_args)
                                              .stdout(process::Stdio::inherit())
                                              .stderr(process::Stdio::inherit())
                                              .status();
/
    // proper error code propogation
    // spaghettified because Command::status() returns - Ok(Exitstatus(Exitstatus(code)))
    // we need to unwrap highest Ok then unwrap the Exitstatus to get to the return code of
    // command executed
    process::exit(match exec_status_result {
        Ok(exec_out) => match exec_out.code() {
            Some(code) => code,
            None       => {
                println!("Process terminated by signal");
                130 // because when I press ctrl-c and echo $? right after, it prints 130
            },
        },
        Err(e) => {
            println!("Error executing command: {}", e);
            2
        },
    });
*/
}

fn print_usage() {
    println!("Usage: ffind [-fdri] [-e=maxdepth] [--debug --help] [expression] [path]");
}

#[derive(Debug)]
struct SortedArgs {
    short_args: Vec<String>,
    long_args: String,
    file_name: String,
    path: Vec<String>,
    exec_args: String,
}

fn make_command(args_in: Vec<String>) -> Result<Vec<String>, io::Error> {
    let args_out = Vec::new();
    let sorted_args = sort_args(args_in)?;
    println!("{:?}", sorted_args);
    Ok(args_out)
}

fn sort_args(args_in: Vec<String>) -> Result<SortedArgs, io::Error> {

    let mut args_in_for_exec = args_in.to_vec();
    let mut short_args: Vec<String> = Vec::new();
    let mut long_args: Vec<String> = Vec::new();
    let mut file_name = String::from("");
    let mut path = Vec::new();
    let mut exec_args = String::from("");

    for (arg_no, arg) in args_in.iter().enumerate() {
        match arg.rsplit_once('-') {
            Some((first, second)) => {
                if first == "-" {
                    long_args.push(second.to_string());
                } else {
                    if second == "exec" {
                        exec_args = args_in_for_exec.drain(arg_no..)
                                                    .collect::<Vec<String>>()
                                                    .join(" ");
                        break;
                    } else {
                        short_args.push(second.to_string());
                    }
                }                
            }
            None => {
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

    let sorted_args = SortedArgs {
        short_args: get_short_args(short_args)?,
        long_args: process_long_args(long_args)?,
        file_name: file_name,
        path: path,
        exec_args: exec_args,
    };

    Ok(sorted_args)
}

fn get_short_args(args_in: Vec<String>) -> Result<Vec<String>, io::Error> {
    let mut ret_short_args: Vec<String> = Vec::new();
    println!("test={:?}", args_in);

    // for args in args_in {
    //     let args_c = args.char_indices();
    //     if args_c.count() > 1 { // short flags or long flags
    //         if args_c.next().unwrap() == (0, '-') {
    //             let second_args_char = args_c.next().unwrap();
    //             if second_args_char == (1, '-') { // second dash so long flag
    //                 continue;
    //             } else {
    //                 ret_short_args.push(second_args_char.1.to_string());
    //                 return Ok(ret_short_args);
    //             }
    //         }
    //     }
    // }
    return Ok(ret_short_args);
}

fn process_long_args(arg_in: Vec<String>) -> Result<String, io::Error> {
    for arg in arg_in {
        match arg.as_str() {
            "help" => {
                print_usage();
                process::exit(0);
            },
            "debug" => {
                // debug logic to be written
            },
            _ => {
                return Err(io::Error::new(io::ErrorKind::Other, format!("flag --{} not recognized", arg)));
            },
        }
    };
    Ok(String::from("long arguments are processed"))
}