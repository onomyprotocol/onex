# Oracle Cloud Ubuntu Firewall Config
IPTABLES_CONFIG=/etc/iptables/rules.v4
if test -f "$IPTABLES_CONFIG"; then
  # Oracle Cloud Ubuntu Firewall Config
  sudo sed -i 's/22 -j ACCEPT/&\n-A INPUT -p tcp -m state --state NEW -m tcp --dport 9091 -j ACCEPT/' $IPTABLES_CONFIG
  sudo sed -i 's/22 -j ACCEPT/&\n-A INPUT -p tcp -m state --state NEW -m tcp --dport 9191 -j ACCEPT/' $IPTABLES_CONFIG
  sudo sed -i 's/22 -j ACCEPT/&\n-A INPUT -p tcp -m state --state NEW -m tcp --dport 1317 -j ACCEPT/' $IPTABLES_CONFIG
  sudo sed -i 's/22 -j ACCEPT/&\n-A INPUT -p tcp -m state --state NEW -m tcp --dport 26657 -j ACCEPT/' $IPTABLES_CONFIG
  sudo iptables-restore < $IPTABLES_CONFIG
fi