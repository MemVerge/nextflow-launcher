# Use Node.js 20 Alpine as the base image
FROM node:20-alpine

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the application
RUN npm run build

# Install serve to run the built application
RUN npm install -g serve

# Expose the port the app runs on
EXPOSE 3000

# Command to run the application
CMD ["serve", "-s", "dist", "-l", "3000"] 