# opensearch-scaling-manager

### Build, Packaging and installation
To install the scaling manager please download the source code using following command:
```bash
git clone https://github.com/maplelabs/opensearch-scaling-manager.git -b release_v0.1_dev
```
Run the following commands to build and install the scaling manager
```bash
cd opensearch-scaling-manager/
# Build the scaling_manager module.
sudo make build
# Package the scaling_manager module and create a tarball.
sudo make pack
# Install the scaling_manager module and create systemd service.
sudo make install
```
To start scaling manager run the following command:
```bash
sudo systemctl start scaling_manager
```
To stop the scaling manager run the following command:
```bash
sudo systemctl stop scaling_manager
```