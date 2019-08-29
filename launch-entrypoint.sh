#!/usr/bin/env bash

box()
{
  local s=("$@") b w
  for l in "${s[@]}"; do
    ((w<${#l})) && { b="$l"; w="${#l}"; }
  done
  tput bold
  tput setaf 3
  echo "    -${b//?/-}-
   | ${b//?/ } |"
  for l in "${s[@]}"; do
    printf '   | %s%*s%s |\n' "$(tput setaf 1)" "-$w" "$l" "$(tput setaf 3)"
  done
  echo "   | ${b//?/ } |
    -${b//?/-}-"
  tput sgr 0
}

show_exit_warning() {
echo "$(tput setaf 2)
||====================================================================||
||//$\\\\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\//$\\\\||
||(100)==================| FEDERAL RESERVE NOTE |================(100)||
||\\\\$//        ~         '------========--------'                \\\\$//||
||<< /        /$\              // ____ \\\\                         \\ >>||
||>>|  12    //L\\\\            // ///..) \\\\         L38036133B   12 |<<||
||<<|        \\\\ //           || <||  >\  ||                        |>>||
||>>|         \\\$/            ||  \$\$ --/  ||        One Hundred     |<<||
||<<|      L38036133B        *\\\\  |\_/  //* series                 |>>||
||>>|  12                     *\\\\/___\_//*   1989                  |<<||
||<<\      Treasurer     ______/Franklin\________     Secretary 12 />>||
||//$\                 ~|UNITED STATES OF AMERICA|~               /$\\\\||
||(100)===================  ONE HUNDRED DOLLARS =================(100)||
||\\\\$//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\\\\$//||
||====================================================================||
$(tput sgr0)"

box "   If you created any infrastructure and did not destroy it " \
    "   you will be accruing charges in your AWS account"
}


trap show_exit_warning EXIT
trap show_exit_warning SIGTERM

bash

